package main

import (
	"flag"
	log "github.com/andriyg76/glog"
	"github.com/andriyg76/scm-backup/git"
	"github.com/andriyg76/scm-backup/ssh"
	"path/filepath"
	"strings"
)

var directoriesStr string
var gitUser, gitEmail string
var sshKey, sshKeyPassword, sshKeyFile string
var trace, debug bool

func init() {
	flag.StringVar(&directoriesStr, "directories", ".", "list of directories to process, split by ','")
	flag.StringVar(&gitUser, "git_user", "", "git actor username")
	flag.StringVar(&gitEmail, "git_email", "", "git actor email")
	flag.StringVar(&git.Username, "git_login", "", "git http/s username")
	flag.StringVar(&git.Password, "git_password", "", "git http/s password")
	flag.StringVar(&sshKey, "ssh_private_key", "", "ssh private key")
	flag.StringVar(&sshKeyFile, "ssh_key_file", "", "ssh private key file")
	flag.StringVar(&sshKeyPassword, "ssh_key_password", "", "ssh private key")
	flag.BoolVar(&trace, "trace", false, "trace logs")
	flag.BoolVar(&debug, "debug", false, "debug logs")
}

func main() {
	flag.Parse()

	var directories []string
	for _, d := range strings.Split(directoriesStr, ",") {
		d = strings.TrimSpace(d)

		if d != "" {
			directories = append(directories, d)
		}
	}

	if trace {
		log.SetLevel(log.TRACE)
	} else if debug {
		log.SetLevel(log.DEBUG)
	}

	log.Debug("Backuping directories: %s", directories)
	abs, _ := filepath.Abs(".")
	log.Debug("Base directory: %s", abs)
	log.Debug("Log level: %s", log.Default())
	if sshKey != "" {
		keyPwd := ""
		if sshKeyPassword != "" {
			keyPwd = "(with keypassword)"
		}
		log.Debug("git auth with ssh key %", keyPwd)
	}
	if git.Username != "" && git.Password != "" {
		log.Debug("git auth with ssh key")
	}

	var agent ssh.SshAgent
	if sshKey != "" {
		var err error
		if err, agent = ssh.CheckSshAgentOrRun(); err != nil {
			log.Fatal("Error starting ssh agent, %s", err)
		}
		defer agent.Stop()
		if err = agent.AddSshKey(sshKey, sshKeyFile, sshKeyPassword); err != nil {
			log.Fatal("Error loading ssh key, %s", err)
		}
	}

	if e := git.CheckUser(gitUser, gitEmail); e != nil {
		log.Fatal("Can't initialise git user parameters %s", e)
	}

	errs := make(map[string]error)
	for _, dir := range directories {
		if e := git.Backup(dir); e != nil {
			errs[dir] = e
		}
	}

	if len(errs) != 0 {
		log.Fatal("Error happened while backing dirs %s", errs)
	}
}
