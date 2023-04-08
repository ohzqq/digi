package tui

import (
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

func ListCollections() []int {
	cols := db.Collections()
	sel := List(cols.RootNames)
	var ids []int
	for _, id := range sel {
		ids = append(ids, cols.RootIDs[id])
	}
	return ids
}

func ListAlbums(al *db.Albums) []int {
	opts := []props.Opt{
		props.Height(10),
		props.ChoiceMap(al.Names),
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

	var ids []int
	for _, id := range l.Selections() {
		ids = append(ids, al.Albums[id].ID)
	}
	return ids
}
