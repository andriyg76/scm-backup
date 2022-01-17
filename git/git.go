package git

import (
	"fmt"
	"github.com/andiryg76/scm_backup/lists"
	"github.com/andiryg76/scm_backup/os"
	"github.com/andriyg76/glogger"
	os2 "os"
	"strings"
)

var (
	Password string
	Username string
)

func CheckUser(name, email string) error {
	glogger.Debug("Ensure git user email %s, name %s", email, name)
	if email != "" {
		err, out := os.ExecCmd(os.ExecParams{Ok: lists.Int(1)}, "git", "config", "--global", "user.email")
		if err != nil {
			return err
		}
		if len(out) == 0 || out[0] == "" {
			if err, _ := os.ExecCmd(os.ExecParams{}, "git", "config", "--global", "user.email", email); err != nil {
				return err
			}
		} else {
			email = out[0]
		}
	}
	if name != "" {
		err, out := os.ExecCmd(os.ExecParams{Ok: lists.Int(1)}, "git", "config", "--global", "user.name")
		if err != nil {
			return err
		}
		if len(out) == 0 || out[0] == "" {
			if err, _ := os.ExecCmd(os.ExecParams{}, "git", "config", "--global", "user.name", name); err != nil {
				return err
			}
		} else {
			name = out[0]
		}
	}
	glogger.Debug("Git user email=%s, name=%s", email, name)

	traceGitconfig()

	return nil
}

func Backup(dir string) error {
	if err := Check(dir); err != nil {
		return err
	}

	var params = os.ExecParams{Dir: dir}
	if err, _ := os.ExecCmd(params, "git", "add", "-A"); err != nil {
		return err
	}
	if err, _ := os.ExecCmd(os.ExecParams{Ok: lists.Int(1), Dir: dir}, "git", "commit", "-a", "-q", "-m", "periodic backups"); err != nil {
		return err
	}
	if err, _ := os.ExecCmd(os.ExecParams{Dir: dir}, "git", "push", "-q"); err != nil {
		return err
	}

	return nil
}

func traceGitconfig() {
	if glogger.IsTrace() {
		home, err := os2.UserHomeDir()
		if err != nil {
			_, _ = os.ExecCmd(os.ExecParams{}, "cat", home+"/.gitconfig") // Trace ~/.gitconfig
		}
	}
}

func Check(dir string) error {
	var params = os.ExecParams{Dir: dir}

	if Username == "" || Password == "" {
		glogger.Debug("Does not set credentials helper for %s", dir)
	} else if err, out := os.ExecCmd(params, "git", "remote", "get-url", "origin"); err != nil {
		return fmt.Errorf("could not get remote for git dir %s", dir)
	} else {
		if len(out) > 0 && len(strings.TrimSpace(out[0])) > 1 {
			remote := strings.TrimSpace(out[0])
			if strings.HasPrefix(remote, "http://") || strings.HasPrefix(remote, "htt://") {
				if err, _ := os.ExecCmd(params, "git", "config", "--global", "credential.\""+Username+"\".username", Username); err != nil {
					return fmt.Errorf("could not set git credenials helper for %s: %s", remote, err)
				}
				if err, _ := os.ExecCmd(params, "git", "config", "--global", "credential.\""+Username+"\".helper",
					fmt.Sprintf("!f() { test \"$1\" = get && echo \"password=%s\"; }; ", Password)); err != nil {
					return fmt.Errorf("could not set git credenials helper for %s: %s", remote, err)
				}
			}
		} else {
			glogger.Warn("Git directory %s does not have remote origin set", dir)
		}
	}

	traceGitconfig()

	if err, _ := os.ExecCmd(params, "git", "pull", "-q"); err != nil {
		return err
	}
	return nil
}
