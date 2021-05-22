package host

import (
	"log"
	"os"
	"sort"

	"golang.org/x/crypto/ssh"
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

func newServiceIndex(hosts yamlConfig) ServiceIndex {
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

type Host struct {
	Hostname string
	IP       string
	Port     string
	Cmd      string
	Cfg      *ssh.ClientConfig
	Services []string
	State    map[string]string
}

type yamlConfig struct {
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

func commandBuilder(services []string, serviceIndex ServiceIndex) string {
	return ""
}

func GetHosts() ([]*Host, *GlobalConfig) {
	yamlFile := parseYAMLConfig()
	srv := newServiceIndex(*yamlFile)
	globalCfg := newGlobalConfig(srv, yamlFile.Interval)

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

	return ret, globalCfg
}
