package db

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

func FileExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func Images() []Image {
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

	var albums []Image
	for rows.Next() {
		var m Image
		err := rows.StructScan(&m)
		if err != nil {
			panic(err)
		}
		albums = append(albums, m)
	}

	return albums
}

const tagsGt = `Tags.id > 21`
const tagsForImg = `Tags.Id IN (
	SELECT tagid
	FROM ImageTags
	WHERE imageid IN (%s)
)`

func tagsForImgSql(id ...int) string {
	var ids []string
	for _, i := range id {
		ids = append(ids, strconv.Itoa(i))
	}
	//sel := sq.Select("tagid").
	//From("ImageTags").
	//Where(sq.Eq{"imageid": id})
	//sql, args, err := sel.ToSql()
	//if err != nil {
	//log.Fatal(err)
	//}
	//fmt.Println(sql)
	//fmt.Println(args)
	return fmt.Sprintf(tagsForImg, strings.Join(ids, ","))
}

func tagsWhere(id ...int) string {
	if len(id) > 0 {
		return fmt.Sprintf("%s AND %s", tagsForImgSql(id...), tagsGt)
	}
	return tagsGt
}

func Tags(ids ...int) []Tag {
	images.mtx.Lock()
	defer images.mtx.Unlock()

	sel := sq.Select(
		"id",
		"Tags.pid as parent",
		"Tags.name as name",
	).
		From("Tags").
		Where(tagsWhere(ids...))
	stmt, args := toSql(sel)

	rows, err := images.DB.Queryx(stmt, args...)
	if err != nil {
		fmt.Println(stmt)
		log.Fatalf("error %v\n", err)
	}
	defer rows.Close()
	images.DB.Unsafe()

	var albums []Tag
	for rows.Next() {
		var m Tag
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
