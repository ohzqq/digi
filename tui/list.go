package tui

import (
	"fmt"
	"log"

	"github.com/londek/reactea"
	"github.com/ohzqq/digi/db"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/props"
)

func List(items []string) []int {
	opts := []props.Opt{
		props.Height(10),
		props.ChoiceSlice(items),
	}
	prop, err := props.New(opts...)
	if err != nil {
		log.Fatal(err)
	}
	prop.NoLimit()
	l := teacozy.New(
		prop,
		teacozy.WithChoice(),
		teacozy.WithFilter(),
	)

	pro := reactea.NewProgram(l)
	if err := pro.Start(); err != nil {
		panic(err)
	}

	return l.Selections()
}

type Model struct {
	db.Collection
	albums db.Albums
}

func New() *Model {
	m := Model{
		Collection: db.GetCollection(),
	}
	m.albums = m.Albums()
	return &m
}

func Start() {
	m := New()
	sel := List(m.ListAlbums().Names())
	fmt.Println(sel)
}

func ListCollections() []int {
	cols := db.GetCollection()
	sel := List(cols.Albums().Names())
	return sel
}

func ListAlbums(al *db.Albums) []int {
	sel := List(al.Names())
	return sel
}
