package db

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ohzqq/digi"
	"github.com/spf13/viper"
)

type Digikam struct {
	DB   *sqlx.DB
	mtx  sync.Mutex
	Path string
}

var (
	images Digikam
	thumbs Digikam
)

const (
	sqliteOpts   = "?cache=shared&mode=ro"
	sqlitePrefix = `file:`
	metaDB       = `digikam4.db`
)

func Connect() {
	images = Digikam{
		Path: filepath.Join(viper.GetString("db"), metaDB),
	}
	if ok := FileExist(images.Path); !ok {
		log.Fatalf("db not found")
	}

	database, err := sqlx.Open("sqlite3", sqlitePrefix+images.Path+sqliteOpts)
	if err != nil {
		log.Fatalf("database connection %v failed\n", err)
	}
	images.DB = database
}

func FileExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func Images() []digi.Image {
	images.mtx.Lock()
	defer images.mtx.Unlock()

	sel := sq.Select(
		"Images.id",
		"Albums.id as album",
		"AlbumRoots.label as root",
		"AlbumRoots.specificPath ||	Albums.relativePath as path",
		"Images.name",
	).
		From("Images").
		InnerJoin(`Albums ON Albums.id = Images.album`).
		InnerJoin(`AlbumRoots ON AlbumRoots.id = Albums.albumRoot`)
	stmt, args := toSql(sel)

	rows, err := images.DB.Queryx(stmt, args...)
	if err != nil {
		fmt.Println(stmt)
		log.Fatalf("error %v\n", err)
	}
	defer rows.Close()
	images.DB.Unsafe()

	var albums []digi.Image
	for rows.Next() {
		var m digi.Image
		err := rows.StructScan(&m)
		if err != nil {
			panic(err)
		}
		albums = append(albums, m)
	}

	return albums
}

func RootAlbums() []digi.AlbumRoot {
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

func Albums() []digi.Album {
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
	stmt, args := toSql(sel)

	rows, err := images.DB.Queryx(stmt, args...)
	if err != nil {
		fmt.Println(stmt)
		log.Fatalf("error %v\n", err)
	}
	defer rows.Close()
	images.DB.Unsafe()

	var albums []digi.Album
	for rows.Next() {
		var m digi.Album
		err := rows.StructScan(&m)
		if err != nil {
			panic(err)
		}
		albums = append(albums, m)
	}
	return albums
}
func toSql(sel sq.SelectBuilder) (string, []any) {
	stmt, args, err := sel.ToSql()
	if err != nil {
		panic(err)
	}
	return stmt, args
}
