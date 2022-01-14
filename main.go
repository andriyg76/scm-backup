package main

import (
	"flag"
	"github.com/andiryg76/scm_backup/git"
	"github.com/andiryg76/scm_backup/os"
	"github.com/andiryg76/scm_backup/ssh"
	log "github.com/andriyg76/glogger"
	"strings"
)

var baseDir string
var directoriesStr string
var directories []string
var gitUser, gitEmail string
var gitLogin, gitPassword string
var sshKey, sshKeyPassword string

func init() {
	flag.StringVar(&baseDir, "base_directory", ".", "basic directory")
	flag.StringVar(&directoriesStr, "directories", ".", "list of directories to process, split by ','")
	directories = strings.Split(directoriesStr, ",")
	flag.StringVar(&gitUser, "git_user", "git", "git actor username")
	flag.StringVar(&gitEmail, "git_email", "email", "git actor email")
	flag.StringVar(&gitLogin, "git_login", "", "git http/s username")
	flag.StringVar(&gitPassword, "git_password", "", "git http/s password")
	flag.StringVar(&sshKey, "ssh_private_key", "", "ssh private key")
	flag.StringVar(&sshKeyPassword, "ssh_key_password", "", "ssh private key")

	flag.Parse()
}

func main() {
	log.SetLevel(log.TRACE)

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

	git.CheckUser(gitUser, gitEmail)

	os.ExecCmd(os.ExecParams{Dir: baseDir}, "git", "status")
}
