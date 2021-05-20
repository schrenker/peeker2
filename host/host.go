package host

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

type Host struct {
	Hostname string
	IP       string
	Port     string
	Cfg      *ssh.ClientConfig
	Services []string

	Cmd   string
	State map[string]string
}

type yamlConfig struct {
	YamlHosts []struct {
		Hostname string   `yaml:"hostname"`
		Ip       string   `yaml:"ip"`
		Port     string   `yaml:"port"`
		User     string   `yaml:"user"`
		KeyPath  string   `yaml:"key"`
		Services []string `yaml:"services"`
	} `yaml:"hosts"`
}

func parseYAMLConfig() *yamlConfig {
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

	var cfg yamlConfig
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalln(err)
	}
	return &cfg
}

func prepareSSHConfig(user, keyPath string) (*ssh.ClientConfig, error) {
	if keyPath == "no" {
		return &ssh.ClientConfig{
			User:            user,
			Auth:            []ssh.AuthMethod{},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}, nil
	}

	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	return &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

// func commandBuilder(services []string) []string

func GetHosts() []*Host {
	yamlFile := parseYAMLConfig()

	ret := make([]*Host, len(yamlFile.YamlHosts))

	for i := range yamlFile.YamlHosts {
		sshcfg, err := prepareSSHConfig(yamlFile.YamlHosts[i].User, yamlFile.YamlHosts[i].KeyPath)
		fmt.Println(sshcfg)
		if err != nil {
			log.Fatalln(err)
		}

		ret[i] = &Host{
			Hostname: yamlFile.YamlHosts[i].Hostname,
			IP:       yamlFile.YamlHosts[i].Ip,
			Port:     yamlFile.YamlHosts[i].Port,
			Services: yamlFile.YamlHosts[i].Services,
			Cfg:      sshcfg,

			State: make(map[string]string),
		}
	}

	return ret
}
