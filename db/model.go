package db

import (
	"fmt"
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
	Images string
}

type Images struct {
	Img []Image
}

type Tags struct {
	Tags []Tag
}

func GetImagesByAlbum(a ...int) Images {
	sel := selectImages()
	if len(a) > 0 {
		sel = sel.Where(sq.Eq{"Albums.id": a})
	}

	return images.GetImages(sel)
}

func GetTagsByImage(ids ...int) Tags {
	sel := sq.Select(
		"id",
		"Tags.pid as parent",
		"Tags.name as name",
	).
		From("Tags").
		Where(tagsWhere(ids...))
	return images.GetTags(sel)
}

const tagsGt = `Tags.id > 21`
const tagsForImg = `Tags.Id IN (
	SELECT tagid
	FROM ImageTags
	WHERE imageid IN (%s)
)`
const whereImgId = `imageid IN (%s)`

func tagsForImgSql(id ...int) string {
	//sel := sq.Select("tagid").
	//From("ImageTags").
	//Where(sq.Eq{"imageid": id})
	//sql, args, err := sel.ToSql()
	//if err != nil {
	//log.Fatal(err)
	//}
	//fmt.Println(sql)
	//fmt.Println(args)
	return fmt.Sprintf(tagsForImg, joinIDs(id))
}

func whereImageId(id ...int) string {
	return fmt.Sprintf(whereImgId, joinIDs(id))
}

func joinIDs(id []int) string {
	var ids []string
	for _, i := range id {
		ids = append(ids, strconv.Itoa(i))
	}
	return strings.Join(ids, ",")
}

func tagsWhere(id ...int) string {
	if len(id) > 0 {
		return fmt.Sprintf("%s AND %s", tagsForImgSql(id...), tagsGt)
	}
	return tagsGt
}

func selectImages() sq.SelectBuilder {
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

	return sel
}
