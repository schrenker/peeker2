package host

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/schrenker/peeker2/config"
	"golang.org/x/crypto/ssh"
)

type Host struct {
	Hostname string
	IP       string
	Port     string
	Cmd      string
	Cfg      *ssh.ClientConfig
	Disks    config.DiskSlice
	Services config.ServiceSlice
	State    map[string]string
}

func (h Host) executeCmd() ([]byte, error) {
	addr := h.Hostname
	if h.IP != "" {
		addr = h.IP
	}
	client, err := ssh.Dial("tcp", addr+":"+h.Port, h.Cfg)
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

func (h *Host) initialState() {
	h.State["hostname"] = h.Hostname
	for k, v := range h.Disks.GetInitialState() {
		h.State[k] = v
	}
	for k, v := range h.Services.GetInitialState() {
		h.State[k] = v
	}
}

func (h *Host) outToState(out []byte) {
	tmp := strings.Split(string(out), "\n")
	h.State["load"] = tmp[0]
	field := 1

	for i := range config.GlobalCfg.DiskIndex {
		if tmp[field] != "" {
			spl := strings.Split(tmp[field], " ")
			atoi, _ := strconv.Atoi(spl[1][:len(spl[1])-1])
			spl[1] = strconv.Itoa(100 - atoi)
			h.State[config.GlobalCfg.DiskIndex[i]] = fmt.Sprintf("%v (%v%%)", spl[0], spl[1])
		} else {
			h.State[config.GlobalCfg.DiskIndex[i]] = tmp[field]
		}
		field++
	}
	for i := range config.GlobalCfg.ServiceIndex {
		h.State[config.GlobalCfg.ServiceIndex[i]] = tmp[field]
		field++
	}
}

func (h *Host) updateState(wg *sync.WaitGroup) {
	defer wg.Done()
	out, err := h.executeCmd()
	if err != nil {
		return
	}
	h.outToState(out)
}
