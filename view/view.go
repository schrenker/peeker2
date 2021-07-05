package view

import (
	"strings"

	tm "github.com/buger/goterm"
	"github.com/schrenker/peeker2/config"
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

func (v view) padding(hosts map[string]*host.Host) {
	for i := range banner {
		max := len(banner[i])

		for _, j := range config.GlobalCfg.HostIndex {
			ln := len(hosts[j].State[banner[i]])
			if ln > max {
				max = ln
			}
		}

		for k := range v {
			v[k][i] = v[k][i] + strings.Repeat(" ", max-len(v[k][i]))
		}
	}
}

func newView(hosts map[string]*host.Host) view {
	ret := make(view, len(hosts)+1)

	ret[0] = banner
	ind := 1

	for _, i := range config.GlobalCfg.HostIndex {
		tmp := make([]string, len(banner))
		for j := range banner {
			tmp[j] = hosts[i].State[banner[j]]
		}
		ret[ind] = tmp
		ind++
	}

	return ret
}
