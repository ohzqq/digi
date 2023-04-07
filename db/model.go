package db

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type AlbumRoots struct {
	ID           int
	Label        string
	SpecificPath string `db:"specificPath"`
}

type Album struct {
	ID        int
	AlbumRoot string `db:"albumRoot"`
}

type Digikam struct {
	DB   *sqlx.DB
	Path string
}

const (
	sqliteOpts   = "?cache=shared&mode=ro"
	sqlitePrefix = `file:`
	metaDB       = `metadata.db`
)

func Connect() Digikam {
	db := Digikam{
		Path: filepath.Join(viper.GetString("db"), "digikam4.db"),
	}
	if ok := util.FileExist(db.Path); !ok {
		log.Fatalf("db not found")
	}

	database, err := sqlx.Open("sqlite3", sqlitePrefix+db.Path+sqliteOpts)
	if err != nil {
		log.Fatalf("database connection %v failed\n", err)
	}
	db.DB = database
	return db
}

func FileExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
