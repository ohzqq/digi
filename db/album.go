package db

import (
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/ohzqq/digi"
)

func (db Digikam) Root() []digi.AlbumRoot {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	sel := sq.Select(
		"id",
		"AlbumRoots.specificPath as path",
		"AlbumRoots.label as name",
	).
		From("AlbumRoots")
	stmt, args := toSql(sel)

	rows, err := db.DB.Queryx(stmt, args...)
	if err != nil {
		fmt.Println(stmt)
		log.Fatalf("error %v\n", err)
	}
	defer rows.Close()
	db.DB.Unsafe()

	var albums []digi.AlbumRoot
	for rows.Next() {
		var m digi.AlbumRoot
		err := rows.StructScan(&m)
		if err != nil {
			panic(err)
		}
		albums = append(albums, m)
	}

	return albums
}

func (db Digikam) Albums() {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	sel := sq.Select("AlbumRoots.specificPath as base", "AlbumRoots.label as parent", "Albums.id", "Albums.relativePath as path").From("Albums").InnerJoin(`AlbumRoots ON AlbumRoots.id = Albums.albumRoot`)
	stmt, args := toSql(sel)

	rows, err := db.DB.Queryx(stmt, args...)
	if err != nil {
		fmt.Println(stmt)
		log.Fatalf("error %v\n", err)
	}
	defer rows.Close()
	db.DB.Unsafe()

	var albums []digi.Album
	for rows.Next() {
		var m digi.Album
		err := rows.StructScan(&m)
		if err != nil {
			panic(err)
		}
		albums = append(albums, m)
	}

	fmt.Printf("%+V\n", albums)
}
