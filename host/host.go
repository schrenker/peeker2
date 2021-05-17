package host

import "golang.org/x/crypto/ssh"

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
		hostname string   `yaml:hostname`
		port     string   `yaml:port`
		user     string   `yaml:user`
		keyPath  string   `yaml:key`
		services []string `yaml:services`
	} `yaml:hosts`
}

func parseYAMLConfig(...string) (*yamlConfig, error) { return nil, nil }

func prepareSSHConfig(user, keyPath string) (*ssh.ClientConfig, error) { return nil, nil }

func GetHosts() ([]*Host, error) { return nil, nil }
