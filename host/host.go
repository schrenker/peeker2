package host

import (
	"fmt"
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
	Disks    []string
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

func commandBuilder(services []string, serviceIndex config.Index, diskIndex config.Index) string {
	cmd := ""
	for i := range serviceIndex {
		if stringInSlice(serviceIndex[i], services) {
			cmd += fmt.Sprintf("systemctl is-active %v;", serviceIndex[i])
		} else {
			cmd += "echo;"
		}
	}
	cmd += "cat /proc/loadavg | awk '{$1 $2 $3}';"
	return cmd
}

func stringInSlice(str string, slice []string) bool {
	for i := range slice {
		if str == slice[i] {
			return true
		}
	}
	return false
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
			Disks:    yamlFile.YamlHosts[i].Disks,
			Cmd:      commandBuilder(yamlFile.YamlHosts[i].Services, globalCfg.ServiceIndex, globalCfg.DiskIndex),
			Cfg:      sshcfg,
			State:    make(map[string]string),
		}
	}

	return ret
}
