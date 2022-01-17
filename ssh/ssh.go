package ssh

import (
	"github.com/andiryg76/scm_backup/lists"
	"github.com/andiryg76/scm_backup/os"
	"github.com/andriyg76/glogger"
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
		glogger.Debug("Ssh agent running socket=%s pid=%s", sock, pid)
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
		glogger.Trace("Skipping to stop agent")
		return
	}
	os.ExecCmd(os.ExecParams{Env: a.env}, "ssh-agent", "-k")
}

func (a SshAgent) AddSshKey(key string, pw string) error {
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