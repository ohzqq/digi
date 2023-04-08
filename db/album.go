package db

import (
	"path/filepath"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

type Collection struct {
	RootIDs   []int
	RootNames []string
	Roots     Albums
	Albums
}

type Root struct {
	ID         int
	Path       string `db:"path"`
	Name       string
	Albums     []Album
	AlbumNames []map[string]string
}

type Album struct {
	ID        int
	AlbumRoot int `db:"-"`
	Parent    string
	Name      string
	Path      string `db:"path"`
	Label     string `db:"label"`
	Base      string `db:"base"`
}

type Albums struct {
	Albums []Album
	Names  []map[string]string
}

func Collections() Collection {
	sel := sq.Select(
		"id",
		"AlbumRoots.specificPath as path",
		"AlbumRoots.label as name",
	).
		From("AlbumRoots")
	return images.GetCollections(sel)
}

func (r *Root) ListAlbums() *Root {
	sel := selectAlbums()
	sel = sel.Where(sq.Eq{"albumRoot": r.ID})

	albums, names := images.GetAlbums(sel)
	r.Albums = albums
	r.AlbumNames = names
	return r
}

func (a Albums) Children() Collection {
	var col Collection
	for i, al := range a.Albums {
		if _, ok := a.Names[i]["/"]; ok {
			col.Roots.Albums = append(col.Roots.Albums, al)
			col.Roots.Names = append(col.Roots.Names, a.Names[i])
		}
	}

	for _, r := range col.Roots.Names {
		var root string
		if n, ok := r["/"]; ok {
			root = n
		}
		for i, al := range a.Albums {
			if _, ok := a.Names[i]["/"]; !ok {
				d := strings.TrimPrefix(al.Path, "/"+root)
				col.Albums.Albums = append(col.Albums.Albums, al)
				col.Names = append(col.Names, map[string]string{d: filepath.Base(al.Path)})
			}
		}
	}
	return col
}

func GetAlbumsByRoot(ids ...int) *Albums {
	r := new(Albums)

	sel := selectAlbums()
	sel = sel.Where(sq.Eq{"albumRoot": ids})

	albums, names := images.GetAlbums(sel)
	r.Albums = albums
	r.Names = names
	return r
}

func GetAlbumsById(ids ...int) Albums {
	sel := selectAlbums()
	if len(ids) > 0 {
		sel = sel.Where(sq.Eq{"Albums.id": ids})
	}
	albums, names := images.GetAlbums(sel)
	return Albums{
		Albums: albums,
		Names:  names,
	}
}

func (a Root) Images() Images {
	var ids []int
	for _, a := range a.Albums {
		ids = append(ids, a.ID)
	}
	return GetImagesByAlbum(ids...)
}

func (a Root) Tags() Tags {
	var ids []int
	for _, img := range a.Images() {
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
		"Albums.relativePath as path",
	).
		From("Albums").
		InnerJoin(`AlbumRoots ON AlbumRoots.id = Albums.albumRoot`)
	return sel
}
