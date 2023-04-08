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
	var items []string
	for _, c := range cols {
		items = append(items, c.Name)
	}
	sel := List(items)
	var ids []int
	for _, id := range sel {
		ids = append(ids, cols[id].ID)
	}
	return ids
}
