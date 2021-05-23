package view

import (
	"time"

	tm "github.com/buger/goterm"
	"github.com/schrenker/peeker2/config"
	"github.com/schrenker/peeker2/host"
)

func generateBanner(disks, services config.Index) []string {
	banner := make([]string, 3+len(disks)+len(services))
	banner[0], banner[1], banner[2] = "1m", "5m", "15m"
	ind := 3
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

func generateView(hosts []*host.Host, banner []string) [][]string { return nil }

func display(view [][]string) {
	tm.Clear()
	tm.MoveCursor(1, 1)
	tm.Flush()
}

func Render(hosts []*host.Host, globalCfg config.GlobalConfig) {
	host.UpdateStatusAll(hosts, globalCfg.DiskIndex, globalCfg.ServiceIndex)
	banner := generateBanner(globalCfg.DiskIndex, globalCfg.ServiceIndex)

	for {
		view := generateView(hosts, banner)
		display(view)
		time.Sleep(time.Duration(globalCfg.Interval) * time.Second)
	}

}
