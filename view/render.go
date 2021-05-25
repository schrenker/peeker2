package view

import (
	"time"

	"github.com/schrenker/peeker2/config"
	"github.com/schrenker/peeker2/host"
)

func generateBanner(disks, services config.Index) []string {
	banner := make([]string, 2+len(disks)+len(services))
	banner[0] = "hostname"
	banner[1] = "load"
	ind := 2
	for i := 0; i < len(disks); i++ {
		banner[ind] = disks[i]
		ind++
	}
	for i := 0; i < len(services); i++ {
		banner[ind] = services[i]
		ind++
	}
	return banner
}

func Render(hosts []*host.Host, globalCfg config.GlobalConfig) {
	host.UpdateStatusAll(hosts, globalCfg.DiskIndex, globalCfg.ServiceIndex)
	banner := generateBanner(globalCfg.DiskIndex, globalCfg.ServiceIndex)

	for {
		view := newView(hosts, banner)
		view.padding(hosts, banner)
		view.display()
		time.Sleep(time.Duration(globalCfg.Interval) * time.Second)
	}

}
