package db

import (
	"strings"

	sq "github.com/Masterminds/squirrel"
)

type Collection struct {
	albums Albums
}

type Album struct {
	ID     int
	Depth  int
	Parent string
	Name   string
	path   []string
	base   []string
	Dir    string
	Path   string `db:"path"`
	Base   string `db:"base"`
}

type Albums []Album

func GetCollection() Collection {
	sel := selectAlbums()
	return images.GetCollection(sel)
}

func (c Collection) Albums() Albums {
	var albums Albums
	for _, a := range c.albums {
		if a.Depth == 0 {
			albums = append(albums, a)
		}
	}
	return albums
}

func (c Collection) ListAlbums() Albums {
	return c.albums
}

func (albums Albums) Names() []string {
	var names []string
	for _, a := range albums {
		names = append(names, strings.Repeat(" ", a.Depth)+a.Name)
	}
	return names
}

func (c Collection) OpenNode(album Album) Albums {
	var albums Albums
	for _, a := range c.albums {
		if a.Base == album.Base {
			if a.Depth == album.Depth+1 {
				albums = append(albums, a)
			}
		}
	}
	return albums
}

func GetAlbumsByRoot(ids ...int) Albums {
	sel := selectAlbums()
	sel = sel.Where(sq.Eq{"albumRoot": ids})
	return images.GetAlbums(sel)
}

func GetAlbumsById(ids ...int) Albums {
	sel := selectAlbums()
	if len(ids) > 0 {
		sel = sel.Where(sq.Eq{"Albums.id": ids})
	}
	return images.GetAlbums(sel)
}

func (albums Albums) Images() Images {
	var ids []int
	for _, a := range albums {
		ids = append(ids, a.ID)
	}
	return GetImagesByAlbum(ids...)
}

func (albums Albums) Tags() Tags {
	var ids []int
	for _, img := range albums.Images() {
		ids = append(ids, img.ID)
	}
	//return GetTagsByImage(ids...)
	return groupImagesByTag(ids...)
}

func groupImagesByTag(ids ...int) Tags {
	sel := sq.Select(
		"Tags.name",
		"GROUP_CONCAT(imageid) as images",
	).
		From("ImageTags").
		InnerJoin("Tags ON tagid = Tags.id").
		Where(whereImageId(ids...) + ` AND ` + tagsGt).
		GroupBy("tagid")
	return images.GetTags(sel)
}

func selectAlbums() sq.SelectBuilder {
	sel := sq.Select(
		"AlbumRoots.specificPath as base",
		"AlbumRoots.label as parent",
		"Albums.id",
		"AlbumRoots.specificPath || Albums.relativePath as path",
	).
		From("Albums").
		InnerJoin(`AlbumRoots ON AlbumRoots.id = Albums.albumRoot`)
	return sel
}
