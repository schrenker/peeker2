package config

import (
	"log"
	"os"
	"sort"

	"gopkg.in/yaml.v2"
)

type GlobalConfig struct {
	ServiceIndex ServiceIndex
	Interval     int
}

func newGlobalConfig(srv ServiceIndex, interval int) *GlobalConfig {
	return &GlobalConfig{
		ServiceIndex: srv,
		Interval:     interval,
	}
}

type ServiceIndex []string

func newServiceIndex(hosts YamlConfig) ServiceIndex {
	var ret ServiceIndex
	amounts := make(map[string]int)

	for i := range hosts.YamlHosts {
		for j := range hosts.YamlHosts[i].Services {
			amounts[hosts.YamlHosts[i].Services[j]]++
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
	Interval  int `yaml:"interval"`
	YamlHosts []struct {
		Hostname string   `yaml:"hostname"`
		Ip       string   `yaml:"ip"`
		Port     string   `yaml:"port"`
		User     string   `yaml:"user"`
		KeyPath  string   `yaml:"key"`
		Services []string `yaml:"services"`
	} `yaml:"hosts"`
}

func parseYAMLConfig() *YamlConfig {
	paths := []string{"./testdata/cfg.yaml", "./cfg.yaml"}
	var yamlFile []byte
	var err error

	for i := range paths {
		yamlFile, err = os.ReadFile(paths[i])
		if yamlFile != nil {
			break
		} else {
			continue
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
	srv := newServiceIndex(*yamlFile)
	globalCfg := newGlobalConfig(srv, yamlFile.Interval)
	return yamlFile, globalCfg
}
