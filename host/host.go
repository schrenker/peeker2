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

func (h Host) executeCmd() ([]byte, error) {
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

func (h *Host) outToState(out []byte, disks, services config.Index) {
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

func (h *Host) updateState(disks, services config.Index) {
	out, err := h.executeCmd()
	if err != nil {
		return
	}
	h.outToState(out, disks, services)
}
