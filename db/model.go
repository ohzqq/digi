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
	"github.com/spf13/viper"
)

type Image struct {
	ID    int
	Album int `db:"album"`
	Root  string
	Path  string `db:"path"`
	Name  string `db:"name"`
}

type Digikam struct {
	DB   *sqlx.DB
	mtx  sync.Mutex
	Path string
}

var (
	images *sqlx.DB
	thumbs *sqlx.DB
)

const (
	sqliteOpts   = "?cache=shared&mode=ro"
	sqlitePrefix = `file:`
	metaDB       = `digikam4.db`
)

func Connect() Digikam {
	db := Digikam{
		Path: filepath.Join(viper.GetString("db"), metaDB),
	}
	if ok := FileExist(db.Path); !ok {
		log.Fatalf("db not found")
	}

	database, err := sqlx.Open("sqlite3", sqlitePrefix+db.Path+sqliteOpts)
	if err != nil {
		log.Fatalf("database connection %v failed\n", err)
	}
	images = database
	db.DB = database
	return db
}

func FileExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func (db Digikam) Images() {
	db.mtx.Lock()
	defer db.mtx.Unlock()

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

	rows, err := db.DB.Queryx(stmt, args...)
	if err != nil {
		fmt.Println(stmt)
		log.Fatalf("error %v\n", err)
	}
	defer rows.Close()
	db.DB.Unsafe()

	var albums []Image
	for rows.Next() {
		var m Image
		err := rows.StructScan(&m)
		if err != nil {
			panic(err)
		}
		albums = append(albums, m)
	}

	fmt.Printf("%+V\n", albums)
}

func toSql(sel sq.SelectBuilder) (string, []any) {
	stmt, args, err := sel.ToSql()
	if err != nil {
		panic(err)
	}
	return stmt, args
}
