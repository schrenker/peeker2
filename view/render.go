package view

import (
	"time"

	"github.com/schrenker/peeker2/config"
	"github.com/schrenker/peeker2/host"
)

var banner []string

func generateBanner(disks, services config.Index) []string {
	banner := make([]string, 2+len(disks)+len(services))
	banner[0] = "hostname"
	banner[1] = "load"
	field := 2
	for i := 0; i < len(disks); i++ {
		banner[field] = disks[i]
		field++
	}
	for i := 0; i < len(services); i++ {
		banner[field] = services[i]
		field++
	}
	return banner
}

func Render(hosts map[string]*host.Host) {
	for {
		host.UpdateStatusAll(hosts)
		banner = generateBanner(config.GlobalCfg.DiskIndex, config.GlobalCfg.ServiceIndex)
		view := newView(hosts)
		view.padding(hosts)
		view.display()
		time.Sleep(time.Duration(config.GlobalCfg.Interval) * time.Second)
	}
}
