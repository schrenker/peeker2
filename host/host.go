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
	Cfg      *ssh.ClientConfig
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

func GetHosts() []*Host {
	yamlFile := parseYAMLConfig()

	ret := make([]*Host, len(yamlFile.yamlHosts))

	for i := range yamlFile.yamlHosts {
		sshcfg, err := prepareSSHConfig(yamlFile.yamlHosts[i].user, yamlFile.yamlHosts[i].keyPath)
		if err != nil {
			log.Fatalln(err)
		}

		ret[i] = &Host{
			Hostname: yamlFile.yamlHosts[i].hostname,
			IP:       yamlFile.yamlHosts[i].ip,
			Port:     yamlFile.yamlHosts[i].port,
			Services: yamlFile.yamlHosts[i].services,
			Cfg:      sshcfg,
		}
	}

	return ret
}
