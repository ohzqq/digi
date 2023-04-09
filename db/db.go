package db

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
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

func (d Digikam) GetAlbums(sel sq.SelectBuilder) Albums {
	images.mtx.Lock()
	defer images.mtx.Unlock()

	sel = sel.OrderBy("parent")

	stmt, args := toSql(sel)
	//fmt.Println(stmt)

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
		m.base = strings.Split(strings.Trim(m.Base, "/"), "/")
		m.path = strings.Split(strings.Trim(m.Path, "/"), "/")
		m.Depth = len(m.path) - len(m.base)
		m.Name = filepath.Base(m.Path)
		if m.Depth == 0 {
			m.Name = m.Parent
		}
		albums = append(albums, m)
	}
	return albums
}

func (d Digikam) GetCollection(sel sq.SelectBuilder) Collection {
	images.mtx.Lock()
	defer images.mtx.Unlock()

	sel = sel.OrderBy("parent")

	stmt, args := toSql(sel)
	//fmt.Println(stmt)

	rows, err := images.DB.Queryx(stmt, args...)
	if err != nil {
		fmt.Println(stmt)
		log.Fatalf("error %v\n", err)
	}
	defer rows.Close()

	var col Collection
	for rows.Next() {
		var m Album
		err := rows.StructScan(&m)
		if err != nil {
			panic(err)
		}
		m.base = strings.Split(strings.Trim(m.Base, "/"), "/")
		m.path = strings.Split(strings.Trim(m.Path, "/"), "/")
		m.Depth = len(m.path) - len(m.base)
		m.Name = filepath.Base(m.Path)
		if m.Depth == 0 {
			m.Name = m.Parent
		}
		col.albums = append(col.albums, m)
	}
	return col
}

func (d Digikam) GetImages(sel sq.SelectBuilder) Images {
	images.mtx.Lock()
	defer images.mtx.Unlock()

	stmt, args := toSql(sel)

	rows, err := images.DB.Queryx(stmt, args...)
	if err != nil {
		fmt.Println(stmt)
		log.Fatalf("error %v\n", err)
	}
	defer rows.Close()
	images.DB.Unsafe()

	var imgs Images
	for rows.Next() {
		var m Image
		err := rows.StructScan(&m)
		if err != nil {
			panic(err)
		}
		imgs = append(imgs, m)
	}

	return imgs
}

func (d Digikam) GetTags(sel sq.SelectBuilder) Tags {
	images.mtx.Lock()
	defer images.mtx.Unlock()

	stmt, args := toSql(sel)
	//fmt.Println(stmt)

	rows, err := images.DB.Queryx(stmt, args...)
	if err != nil {
		fmt.Println(stmt)
		log.Fatalf("error %v\n", err)
	}
	defer rows.Close()
	images.DB.Unsafe()

	var tags Tags
	for rows.Next() {
		var m Tag
		err := rows.StructScan(&m)
		if err != nil {
			panic(err)
		}
		tags = append(tags, m)
	}

	return tags
}

func FileExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func toSql(sel sq.SelectBuilder) (string, []any) {
	stmt, args, err := sel.ToSql()
	if err != nil {
		panic(err)
	}
	return stmt, args
}
