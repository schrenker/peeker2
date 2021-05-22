package host

import (
	"log"
	"os"

	"github.com/schrenker/peeker2/config"
	"golang.org/x/crypto/ssh"
)

type Host struct {
	Hostname string
	IP       string
	Port     string
	Cmd      string
	Cfg      *ssh.ClientConfig
	Services []string
	State    map[string]string
}

func prepareSSHConfig(user, keyPath string) (*ssh.ClientConfig, error) {
	if keyPath == "" {
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

func commandBuilder(services []string, serviceIndex config.ServiceIndex) string {
	return ""
}

func GetHosts(yamlFile config.YamlConfig, globalCfg config.GlobalConfig) []*Host {
	ret := make([]*Host, len(yamlFile.YamlHosts))

	for i := range yamlFile.YamlHosts {
		sshcfg, err := prepareSSHConfig(yamlFile.YamlHosts[i].User, yamlFile.YamlHosts[i].KeyPath)
		if err != nil {
			log.Fatalln(err)
		}

		ret[i] = &Host{
			Hostname: yamlFile.YamlHosts[i].Hostname,
			IP:       yamlFile.YamlHosts[i].Ip,
			Port:     yamlFile.YamlHosts[i].Port,
			Services: yamlFile.YamlHosts[i].Services,
			Cmd:      commandBuilder(yamlFile.YamlHosts[i].Services, globalCfg.ServiceIndex),
			Cfg:      sshcfg,
			State:    make(map[string]string),
		}
	}

	return ret
}
