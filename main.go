package main

import (
	"flag"
	log "github.com/andriyg76/glogger"
	"github.com/andriyg76/scm-backup/git"
	"github.com/andriyg76/scm-backup/ssh"
	"strings"
)

var directoriesStr string
var directories []string
var gitUser, gitEmail string
var sshKey, sshKeyPassword string
var trace, debug bool

func init() {
	flag.StringVar(&directoriesStr, "directories", ".", "list of directories to process, split by ','")
	for _, d := range strings.Split(directoriesStr, ",") {
		d = strings.TrimSpace(d)

		if d != "" {
			directories = append(directories, d)
		}
	}
	flag.StringVar(&gitUser, "git_user", "git", "git actor username")
	flag.StringVar(&gitEmail, "git_email", "email", "git actor email")
	flag.StringVar(&git.Username, "git_login", "", "git http/s username")
	flag.StringVar(&git.Password, "git_password", "", "git http/s password")
	flag.StringVar(&sshKey, "ssh_private_key", "", "ssh private key")
	flag.StringVar(&sshKeyPassword, "ssh_key_password", "", "ssh private key")
	flag.BoolVar(&trace, "trace", false, "trace logs")
	flag.BoolVar(&debug, "debug", false, "debug logs")
	flag.Parse()
}

func main() {
	if trace {
		log.SetLevel(log.TRACE)
	} else if debug {
		log.SetLevel(log.DEBUG)
	}

	var agent ssh.SshAgent
	if sshKey != "" {
		var err error
		if err, agent = ssh.CheckSshAgentOrRun(); err != nil {
			log.Fatal("Error starting ssh agent, %s", err)
		}
		defer agent.Stop()
		if err = agent.AddSshKey(sshKey, sshKeyPassword); err != nil {
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
