package db

import (
	sq "github.com/Masterminds/squirrel"
)

type Collection struct {
	Parent string
	Name   string
	albums Albums
	Names  []string
}

type Root struct {
	ID   int
	Path string `db:"path"`
	Name string
}

type Roots struct {
	IDs   []int
	Names []string
	Paths []string
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

func Collections() Collection {
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

func (albums Albums) Names() []string {
	var names []string
	for _, a := range albums {
		names = append(names, a.Name)
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

//func (a Root) Images() Images {
//  var ids []int
//  for _, a := range a.Albums {
//    ids = append(ids, a.ID)
//  }
//  return GetImagesByAlbum(ids...)
//}

//func (a Root) Tags() Tags {
//  var ids []int
//  for _, img := range a.Images() {
//    ids = append(ids, img.ID)
//  }
//  //return GetTagsByImage(ids...)
//  return groupImagesByTag(ids...)
//}

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
