package db

import (
	"fmt"
	"log"
	"path/filepath"

	sq "github.com/Masterminds/squirrel"
)

type Collection struct {
	ID   int
	Path string `db:"path"`
	Name string
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
	images.mtx.Lock()
	defer images.mtx.Unlock()

	sel := sq.Select(
		"id",
		"AlbumRoots.specificPath as path",
		"AlbumRoots.label as name",
	).
		From("AlbumRoots")
	stmt, args := toSql(sel)

	rows, err := images.DB.Queryx(stmt, args...)
	if err != nil {
		fmt.Println(stmt)
		log.Fatalf("error %v\n", err)
	}
	defer rows.Close()
	images.DB.Unsafe()

	var albums []Collection
	for rows.Next() {
		var m Collection
		err := rows.StructScan(&m)
		if err != nil {
			panic(err)
		}
		albums = append(albums, m)
	}

	return albums
}

func (r *Collection) Albums() Albums {
	images.mtx.Lock()
	defer images.mtx.Unlock()

	sel := selectAlbums()
	sel = sel.Where(sq.Eq{"albumRoot": r.ID})

	stmt, args := toSql(sel)
	fmt.Println(stmt)

	rows, err := images.DB.Queryx(stmt, args...)
	if err != nil {
		fmt.Println(stmt)
		log.Fatalf("error %v\n", err)
	}
	defer rows.Close()

	var albums Albums
	for rows.Next() {
		var m Album
		err := rows.StructScan(&m)
		if err != nil {
			panic(err)
		}
		albums.Names = append(albums.Names, filepath.Base(m.Path))
		albums.Albums = append(albums.Albums, m)
	}
	return albums
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

func (a Albums) Images() Images {
	var ids []int
	for _, a := range a.Albums {
		ids = append(ids, a.ID)
	}
	return GetImages(ids...)
}

func (a Albums) Tags() Tags {
	var ids []int
	for _, img := range a.Images().Img {
		ids = append(ids, img.ID)
	}
	return GetTags(ids...)
}

func GetAlbums(ids ...int) Albums {
	images.mtx.Lock()
	defer images.mtx.Unlock()

	sel := selectAlbums()
	if len(ids) > 0 {
		sel = sel.Where(sq.Eq{"Albums.id": ids})
	}
	return images.GetAlbums(sel)
}
