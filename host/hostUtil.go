package host

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/schrenker/peeker2/config"
	"golang.org/x/crypto/ssh"
)

var Embedded embed.FS

func prepareSSHConfig(user, keyPath string) (*ssh.ClientConfig, error) {
	if keyPath == "" {
		return &ssh.ClientConfig{
			User:            user,
			Auth:            []ssh.AuthMethod{},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}, nil
	}

	var key []byte
	var err error

	if strings.Contains(keyPath, "embed") {
		key, err = Embedded.ReadFile(keyPath)
		if err != nil {
			return nil, err
		}
	} else {
		key, err = os.ReadFile(keyPath)
		if err != nil {
			return nil, err
		}
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

func commandBuilder(disks, services config.Observable, diskIndex, serviceIndex config.Index) string {
	cmd := "cat /proc/loadavg | awk '{print $1\" \"$2\" \"$3}';"
	for i := range diskIndex {
		if stringInSlice(diskIndex[i], disks.GetNames()) {
			cmd += fmt.Sprintf("df -hBG | grep -w %v | awk '{print $4\" \"$5}';", diskIndex[i])
		} else {
			cmd += "echo;"
		}
	}
	for i := range serviceIndex {
		if stringInSlice(serviceIndex[i], services.GetNames()) {
			cmd += fmt.Sprintf("systemctl is-active %v;", serviceIndex[i])
		} else {
			cmd += "echo;"
		}
	}
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
	ret := make([]*Host, len(yamlFile.Hosts))

	for i := range yamlFile.Hosts {
		sshcfg, err := prepareSSHConfig(yamlFile.Hosts[i].User, yamlFile.Hosts[i].KeyPath)
		if err != nil {
			log.Fatalln(err)
		}

		ret[i] = &Host{
			Hostname: yamlFile.Hosts[i].Hostname,
			IP:       yamlFile.Hosts[i].Ip,
			Port:     yamlFile.Hosts[i].Port,
			Services: yamlFile.Hosts[i].Services,
			Disks:    yamlFile.Hosts[i].Disks,
			Cmd: commandBuilder(
				yamlFile.Hosts[i].Disks,
				yamlFile.Hosts[i].Services,
				globalCfg.DiskIndex,
				globalCfg.ServiceIndex,
			),
			Cfg:   sshcfg,
			State: make(map[string]string),
		}
		ret[i].initialState()
	}

	return ret
}

func UpdateStatusAll(hosts []*Host, disks, services config.Index) {
	var wg sync.WaitGroup

	for i := range hosts {
		wg.Add(1)
		go hosts[i].updateState(disks, services, &wg)
	}

	wg.Wait()
}
