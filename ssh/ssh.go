package ssh

import (
	"bufio"
	"fmt"
	glog "github.com/andriyg76/glog"
	"github.com/andriyg76/scm-backup/lists"
	"github.com/andriyg76/scm-backup/os"
	"io"
	os2 "os"
	"strings"
)

type SshAgent struct {
	socket string
	env    []string
}

func CheckSshAgentOrRun() (error, SshAgent) {
	agent := SshAgent{}

	sock, pid := os2.Getenv("SSH_AUTH_SOCK"), os2.Getenv("SSH_AGENT_PID")
	if sock != "" && pid != "" {
		glog.Debug("Ssh agent running socket=%s pid=%s", sock, pid)
		agent.socket = sock
		return nil, agent
	}

	err, out := os.ExecCmd(os.ExecParams{}, "ssh-agent")
	if err != nil {
		return err, agent
	}
	return nil, getAgetEnv(out)
}

func (a SshAgent) Stop() {
	if len(a.env) == 0 {
		glog.Trace("Skipping to stop agent")
		return
	}
	os.ExecCmd(os.ExecParams{Env: a.env}, "ssh-agent", "-k")
}

func (a SshAgent) AddSshKey(key string, fileName string, pw string) error {
	if key != "" {
		if err := a.addKey(key, pw); err != nil {
			return err
		}
	}
	if fileName != "" {
		if err := a.addKeyFile(fileName, pw); err != nil {
			return err
		}
	}

	return nil
}

func (a SshAgent) addKey(key string, pw string) error {
	var stdin string
	var args []string
	if pw != "" { // We are removing passwork from private key and adding it from temporary file then
		var err error
		var file *os2.File
		file, err = os2.CreateTemp(os2.TempDir(), "key")
		if err != nil {
			return err
		}
		defer func() {
			os2.Remove(file.Name())
		}()
		file.Chmod(0600)
		file.WriteString(key)
		file.Close()
		if err, _ = os.ExecCmd(os.ExecParams{}, "ssh-keygen", "-p", "-P", pw, "-N", "", "-f", file.Name()); err != nil {
			return err
		}

		stdin = ""
		args = lists.String(file.Name())
	} else {
		stdin = key
		args = lists.String("-")
	}
	err, _ := os.ExecCmd(os.ExecParams{Stdin: stdin, Env: a.env}, "ssh-add", args...)
	return err
}

func (a SshAgent) addKeyFile(fileName string, pw string) error {

	origin, err2 := os2.OpenFile(fileName, os2.O_RDONLY, 0)
	if err2 != nil {
		return fmt.Errorf("can't load original key file %s", err2)
	}

	// We are  adding key from temporary file, as original could have invalid permissions
	var err error
	var file *os2.File
	file, err = os2.CreateTemp(os2.TempDir(), "key")
	if err != nil {
		return err
	}
	defer func() {
		os2.Remove(file.Name())
	}()
	file.Chmod(0600)

	for {
		written, err3 := io.Copy(bufio.NewWriter(file), bufio.NewReader(origin))
		if err3 == io.EOF || written == 0 {
			break
		} else if err3 != nil {
			return fmt.Errorf("can't copy original key file %s", err3)
		}
	}
	origin.Close()
	file.Close()

	// We are removing passwork from private key
	if pw != "" {
		if err, _ = os.ExecCmd(os.ExecParams{}, "ssh-keygen", "-p", "-P", pw, "-N", "", "-f", file.Name()); err != nil {
			return err
		}
	}

	args := lists.String(file.Name())
	err, _ = os.ExecCmd(os.ExecParams{Stdin: "", Env: a.env}, "ssh-add", args...)
	return err
}

func getAgetEnv(out []string) SshAgent {
	var env []string
	var socket string
	for _, line := range out {
		for _, statement := range strings.Split(line, ";") {
			if strings.HasPrefix(statement, "SSH_AUTH_SOCK=") {
				env = append(env, strings.TrimSpace(statement))
				socket = strings.TrimSpace(strings.TrimPrefix(statement, "SSH_AUTH_SOCK="))
			} else if strings.HasPrefix(statement, "SSH_AGENT_PID=") {
				env = append(env, strings.TrimSpace(statement))
			}
		}
	}
	return SshAgent{env: env, socket: socket}
}
