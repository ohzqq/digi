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
	return GetAlbums(r.ID)
}

//func (a Albums) Images() []Image {
//}

func GetAlbums(ids ...int) Albums {
	images.mtx.Lock()
	defer images.mtx.Unlock()

	sel := sq.Select(
		"AlbumRoots.specificPath as base",
		"AlbumRoots.label as parent",
		"Albums.id",
		"Albums.relativePath as path",
	).
		From("Albums").
		InnerJoin(`AlbumRoots ON AlbumRoots.id = Albums.albumRoot`)
	if len(ids) > 0 {
		sel = sel.Where(sq.Eq{"albumRoot": ids})
	}
	stmt, args := toSql(sel)

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
