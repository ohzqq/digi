package db

import (
	sq "github.com/Masterminds/squirrel"
)

type Collection struct {
	ID         int
	Path       string `db:"path"`
	Name       string
	Albums     []Album
	AlbumNames []string
}

type Album struct {
	ID        int
	AlbumRoot int `db:"-"`
	Parent    string
	Path      string `db:"path"`
	Label     string `db:"label"`
	Base      string `db:"base"`
}

type Albums struct {
	Albums []Album
	Names  []string
}

func Collections() []Collection {
	sel := sq.Select(
		"id",
		"AlbumRoots.specificPath as path",
		"AlbumRoots.label as name",
	).
		From("AlbumRoots")
	return images.GetCollections(sel)
}

func (r *Collection) ListAlbums() *Collection {
	sel := selectAlbums()
	sel = sel.Where(sq.Eq{"albumRoot": r.ID})

	albums, names := images.GetAlbums(sel)
	r.Albums = albums
	r.AlbumNames = names
	return r
}

func GetAlbumsByRoot(ids ...int) *Collection {
	r := new(Collection)

	sel := selectAlbums()
	sel = sel.Where(sq.Eq{"albumRoot": ids})

	albums, names := images.GetAlbums(sel)
	r.Albums = albums
	r.AlbumNames = names
	return r
}

func GetAlbumsById(ids ...int) Collection {
	sel := selectAlbums()
	if len(ids) > 0 {
		sel = sel.Where(sq.Eq{"Albums.id": ids})
	}
	albums, names := images.GetAlbums(sel)
	return Collection{
		Albums:     albums,
		AlbumNames: names,
	}
}

func (a Collection) Images() Images {
	var ids []int
	for _, a := range a.Albums {
		ids = append(ids, a.ID)
	}
	return GetImagesByAlbum(ids...)
}

func (a Collection) Tags() Tags {
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
