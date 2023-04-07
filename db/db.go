package db

import (
	"errors"
	"log"
	"os"
	"path/filepath"
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

func toSql(sel sq.SelectBuilder) (string, []any) {
	stmt, args, err := sel.ToSql()
	if err != nil {
		panic(err)
	}
	return stmt, args
}
