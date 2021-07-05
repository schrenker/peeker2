package view

import (
	"time"

	"github.com/schrenker/peeker2/config"
	"github.com/schrenker/peeker2/host"
)

// var banner []string

func generateBanner() []string {
	banner := make([]string, 2+len(config.GlobalCfg.DiskIndex)+len(config.GlobalCfg.ServiceIndex))
	banner[0] = "hostname"
	banner[1] = "load"
	field := 2
	for i := 0; i < len(config.GlobalCfg.DiskIndex); i++ {
		banner[field] = config.GlobalCfg.DiskIndex[i]
		field++
	}
	for i := 0; i < len(config.GlobalCfg.ServiceIndex); i++ {
		banner[field] = config.GlobalCfg.ServiceIndex[i]
		field++
	}
	return banner
}

func Render(hosts map[string]*host.Host) {
	for {
		host.UpdateStatusAll(hosts)
		view := newView(hosts)
		view.display()
		time.Sleep(time.Duration(config.GlobalCfg.Interval) * time.Second)
	}
}
