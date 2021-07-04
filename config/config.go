package config

import (
	"embed"
	"log"
	"os"
	"sort"

	"gopkg.in/yaml.v2"
)

var Embedded embed.FS
var ConfigFile string

type GlobalConfig struct {
	ServiceIndex Index
	DiskIndex    Index
	Interval     int
}

func newGlobalConfig(disk, srv Index, interval int) *GlobalConfig {
	return &GlobalConfig{
		DiskIndex:    disk,
		ServiceIndex: srv,
		Interval:     interval,
	}
}

type Index []string

func newIndex(hosts YamlConfig, indexType string) Index {
	var ret Index
	amounts := make(map[string]int)

	for i := range hosts.Hosts {
		switch indexType {
		case "service":
			for j := range hosts.Hosts[i].Services {
				amounts[hosts.Hosts[i].Services[j].Name]++
			}
		case "disk":
			for j := range hosts.Hosts[i].Disks {
				amounts[hosts.Hosts[i].Disks[j].Path]++
			}
		}
	}

	tmp := make([]string, 0, len(amounts))
	for i := range amounts {
		tmp = append(tmp, i)
	}
	sort.Slice(tmp, func(i, j int) bool {
		return amounts[tmp[i]] > amounts[tmp[j]]
	})

	for i := range tmp {
		ret = append(ret, tmp[i])
	}

	return ret
}

type YamlConfig struct {
	Interval int `yaml:"interval"`
	Hosts    []struct {
		Hostname string       `yaml:"hostname"`
		Ip       string       `yaml:"ip"`
		Port     string       `yaml:"port"`
		User     string       `yaml:"user"`
		KeyPath  string       `yaml:"key"`
		Disks    DiskSlice    `yaml:"disks"`
		Services ServiceSlice `yaml:"services"`
	} `yaml:"hosts"`
}

func parseYAMLConfig() *YamlConfig {
	// paths := make([]string, 2, 3)
	// if len(os.Args) > 1 {
	// 	paths = append(paths, os.Args[1])
	// }
	// paths = append(paths, []string{"./testdata/cfg.yaml", "./cfg.yaml"}...)

	paths := []string{ConfigFile, "./cfg.yaml", "./testdata/cfg.yaml"}

	var yamlFile []byte
	var err error

	if yamlFile, err = Embedded.ReadFile("embed/cfg.yaml"); err != nil {
		for i := range paths {
			yamlFile, err = os.ReadFile(paths[i])
			if yamlFile != nil {
				break
			} else {
				continue
			}
		}
	}

	if yamlFile == nil {
		log.Fatalln(err)
	}

	var cfg YamlConfig
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalln(err)
	}
	return &cfg
}

func GetConfig() (*YamlConfig, *GlobalConfig) {
	yamlFile := parseYAMLConfig()
	disk := newIndex(*yamlFile, "disk")
	srv := newIndex(*yamlFile, "service")
	globalCfg := newGlobalConfig(disk, srv, yamlFile.Interval)
	return yamlFile, globalCfg
}
