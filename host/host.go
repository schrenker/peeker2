package host

import (
	"log"
	"os"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

type Host struct {
	Hostname string
	IP       string
	Port     string

	Cfg *ssh.ClientConfig

	Services []string

	Cmd   string
	State map[string]string
}

type yamlConfig struct {
	yamlHosts []struct {
		hostname string   `yaml:"hostname"`
		ip       string   `yaml:"ip"`
		port     string   `yaml:"port"`
		user     string   `yaml:"user"`
		keyPath  string   `yaml:"key"`
		services []string `yaml:"services"`
	} `yaml:"hosts"`
}

func parseYAMLConfig() *yamlConfig {
	paths := []string{"./testdata/cfg.yaml", "./cfg.yaml", os.Args[1]}
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

	var cfg yamlConfig
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalln(err)
	}
	return &cfg
}

// func prepareSSHConfig(user, keyPath string) (*ssh.ClientConfig, error) { return nil, nil }

// func parseHost(services []string) {}

// func GetHosts() ([]*Host, error) {
// 	yaml := parseYAMLConfig()
// }
