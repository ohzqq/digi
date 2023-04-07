package db

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

type Image struct {
	ID    int
	Album int `db:"album"`
	Root  string
	Path  string `db:"path"`
	Name  string `db:"name"`
}

type Tag struct {
	ID     int
	Parent int
	Name   string
}

type Images struct {
	Img []Image
}

type Tags struct {
	Tags []Tag
}

func GetImages(a ...int) Images {
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

	if len(a) > 0 {
		sel = sel.Where(sq.Eq{"Albums.id": a})
	}

	stmt, args := toSql(sel)

	rows, err := images.DB.Queryx(stmt, args...)
	if err != nil {
		fmt.Println(stmt)
		log.Fatalf("error %v\n", err)
	}
	defer rows.Close()
	images.DB.Unsafe()

	var albums Images
	for rows.Next() {
		var m Image
		err := rows.StructScan(&m)
		if err != nil {
			panic(err)
		}
		albums.Img = append(albums.Img, m)
	}

	return albums
}

func GetTags(ids ...int) Tags {
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
	//fmt.Println(stmt)

	rows, err := images.DB.Queryx(stmt, args...)
	if err != nil {
		fmt.Println(stmt)
		log.Fatalf("error %v\n", err)
	}
	defer rows.Close()
	images.DB.Unsafe()

	var albums Tags
	for rows.Next() {
		var m Tag
		err := rows.StructScan(&m)
		if err != nil {
			panic(err)
		}
		albums.Tags = append(albums.Tags, m)
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
