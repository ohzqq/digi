package db

import (
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
)

type AlbumRoot struct {
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

func (db Digikam) Root() []AlbumRoot {
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

	var albums []AlbumRoot
	for rows.Next() {
		var m AlbumRoot
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

	var albums []Album
	for rows.Next() {
		var m Album
		err := rows.StructScan(&m)
		if err != nil {
			panic(err)
		}
		albums = append(albums, m)
	}

	fmt.Printf("%+V\n", albums)
}
