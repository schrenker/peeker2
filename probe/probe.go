package probe

import (
	"bytes"
	"fmt"

	"github.com/schrenker/peeker2/host"
	"golang.org/x/crypto/ssh"
)

func executeCmd(cmd, addr, port string, cfg *ssh.ClientConfig) ([]byte, error) {
	client, err := ssh.Dial("tcp", addr+":"+port, cfg)
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
	session.Run(cmd)

	return stdoutBuf.Bytes(), nil
}

func ProbeHost(host *host.Host) {
	out, _ := executeCmd(host.Cmd, host.Hostname, host.Port, host.Cfg)
	fmt.Println(string(out))
}
