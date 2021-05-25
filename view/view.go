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
		tm.Println(strings.Join(v[i], "  "))
	}
	tm.Flush()
}

func (v view) padding(hosts []*host.Host, banner []string) {
	for i := range banner {
		max := len(banner[i])

		for j := range hosts {
			ln := len(hosts[j].State[banner[i]])
			if ln > max {
				max = ln
			}
			// if i == 0 {
			// 	if len(hosts[j].Hostname) > max {
			// 		max = len(hosts[j].Hostname)
			// 	}
			// } else {
			// 	if len(hosts[j].State[banner[i]]) > max {
			// 		max = len(hosts[j].State[banner[i]])
			// 	}
			// }
		}

		for k := range v {
			v[k][i] = v[k][i] + strings.Repeat(" ", max-len(v[k][i]))
		}
	}
}

func newView(hosts []*host.Host, banner []string) view {
	ret := make(view, len(hosts)+1)

	ret[0] = banner

	for i := range hosts {
		tmp := make([]string, len(banner))
		// tmp[0] = hosts[i].Hostname
		// for j := 1; j < len(banner); j++ {
		for j := range banner {
			tmp[j] = hosts[i].State[banner[j]]
		}
		ret[i+1] = tmp
	}

	return ret
}

// func colorize(row []string) {}
