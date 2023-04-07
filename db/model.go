package db

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
