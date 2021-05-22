package host

import (
	"bytes"

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
