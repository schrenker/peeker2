package host

import (
	"log"
	"os"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

const (
	LOAD = iota
	MEM
	ROOTFS
	HOMEFS
	VARFS
	EXIM
	//more stats and services
)

type Host struct {
	hostname string
	port     string

	cfg *ssh.ClientConfig

	services []int
}

type yamlConfig struct {
	yamlHosts []struct {
		hostname string   `yaml:"hostname"`
		port     string   `yaml:"port"`
		user     string   `yaml:"user"`
		keyPath  string   `yaml:"key"`
		services []string `yaml:"services"`
	} `yaml:"hosts"`
}

func parseYAMLConfig(...string) *yamlConfig {
	yamlFile, err := os.ReadFile("./testdata/cfg.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	var cfg yamlConfig
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalln(err)
	}
	return &cfg
}

// func prepareSSHConfig(user, keyPath string) (*ssh.ClientConfig, error) { return nil, nil }

// func GetHosts() ([]*Host, error) { return nil, nil }
