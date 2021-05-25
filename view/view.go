package view

import (
	"strings"

	tm "github.com/buger/goterm"
	"github.com/schrenker/peeker2/host"
)

type view [][]string

func (v view) display() {
	tm.Clear()
	tm.MoveCursor(1, 1)
	for i := range v {
		tm.Println(strings.Join(v[i], " "))
	}
	tm.Flush()
}

func newView(hosts []*host.Host, banner []string) view {
	ret := make(view, len(hosts)+1)

	ret[0] = banner

	for i := range hosts {
		tmp := make([]string, len(banner))
		for j := range banner {
			tmp[j] = hosts[i].State[banner[j]]
		}
		ret[i+1] = tmp
	}

	return ret
}

// func colorize(row []string) {}
// func padding(view [][]string) {}
