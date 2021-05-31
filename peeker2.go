package main

import (
	"embed"

	"github.com/schrenker/peeker2/config"
	"github.com/schrenker/peeker2/host"
	"github.com/schrenker/peeker2/view"
)

//go:embed embed/*
var embedded embed.FS

func main() {
	config.Embedded = embedded
	host.Embedded = embedded
	yamlFile, globalCfg := config.GetConfig()
	hosts := host.GetHosts(*yamlFile, *globalCfg)
	view.Render(hosts, *globalCfg)
}
