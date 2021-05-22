package host

import (
	"bytes"
	"strings"

	"github.com/schrenker/peeker2/config"
	"golang.org/x/crypto/ssh"
)

type Host struct {
	Hostname string
	IP       string
	Port     string
	Cmd      string
	Cfg      *ssh.ClientConfig
	Disks    []string
	Services []string
	State    map[string]string
}

func (h Host) ExecuteCmd() ([]byte, error) {
	client, err := ssh.Dial("tcp", h.Hostname+":"+h.Port, h.Cfg)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(h.Cmd)

	return stdoutBuf.Bytes(), nil
}

func (h *Host) OutToState(out []byte, disks, services config.Index) {
	tmp := strings.Split(string(out), "\n")
	h.State["load"] = tmp[0]
	field := 1

	for i := range disks {
		h.State[disks[i]] = tmp[field]
		field++
	}
	for i := range services {
		h.State[services[i]] = tmp[field]
		field++
	}
}
