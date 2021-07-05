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
	for i := range v[0] {
		max := len(v[0][i])

		for _, j := range config.GlobalCfg.HostIndex {
			ln := len(hosts[j].State[v[0][i]])
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
	v := make(view, len(hosts)+1)
	v[0] = generateBanner()
	iter := 1

	for _, i := range config.GlobalCfg.HostIndex {
		tmp := make([]string, len(v[0]))
		for j := range v[0] {
			tmp[j] = hosts[i].State[v[0][j]]
		}
		v[iter] = tmp
		iter++
	}

	v.padding(hosts)

	return v
}
