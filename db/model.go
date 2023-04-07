package db

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
