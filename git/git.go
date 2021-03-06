package git

import (
	"fmt"
	glog "github.com/andriyg76/glog"
	"github.com/andriyg76/scm-backup/lists"
	"github.com/andriyg76/scm-backup/os"
	"net/url"
	os2 "os"
	"strings"
)

var (
	Password string
	Username string
	Env      []string
)

func CheckUser(name, email string) error {
	glog.Debug("Ensure git user email %s, name %s", email, name)
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
	glog.Debug("Git user email=%s, name=%s", email, name)

	traceGitconfig()

	return nil
}

func Backup(dir string) error {
	if err := Check(dir); err != nil {
		return err
	}

	if err, _ := os.ExecCmd(os.ExecParams{Dir: dir}, "git", "add", "-A"); err != nil {
		return err
	}
	if err, _ := os.ExecCmd(os.ExecParams{Ok: lists.Int(1), Dir: dir, Env: Env}, "git", "commit", "-a", "-q", "-m", "periodic backups"); err != nil {
		return err
	}
	if err, _ := os.ExecCmd(os.ExecParams{Dir: dir, Env: Env}, "git", "push", "-q"); err != nil {
		return err
	}

	return nil
}

func traceGitconfig() {
	if glog.IsTrace() {
		home, err := os2.UserHomeDir()
		if err == nil {
			_, _ = os.ExecCmd(os.ExecParams{}, "cat", home+"/.gitconfig") // Trace ~/.gitconfig
		} else {
			glog.Error("Can't get user home dir %s", err)
		}
	}
}

func Check(dir string) error {
	if Username == "" || Password == "" {
		glog.Debug("Does not set credentials helper for %s", dir)
	} else if err, out := os.ExecCmd(os.ExecParams{Dir: dir}, "git", "remote", "get-url", "origin"); err != nil {
		return fmt.Errorf("could not get remote for git dir %s", dir)
	} else {
		if len(out) > 0 && len(strings.TrimSpace(out[0])) > 1 {
			remote := strings.TrimSpace(out[0])
			if strings.HasPrefix(remote, "http://") || strings.HasPrefix(remote, "https://") {
				url, err := url.Parse(remote)
				if err != nil {
					return fmt.Errorf("invalid remote url %s: %s", remote, err)
				}
				remote = url.Scheme + "://" + url.Host
				if err, _ := os.ExecCmd(os.ExecParams{}, "git", "config", "--global", "credential."+remote+".username", Username); err != nil {
					return fmt.Errorf("could not set git credenials helper for %s: %s", remote, err)
				}
				if err, _ := os.ExecCmd(os.ExecParams{}, "git", "config", "--global", "credential."+remote+".helper",
					fmt.Sprintf("!f() { test \"$1\" = get && echo \"password=%s\"; }; f", Password)); err != nil {
					return fmt.Errorf("could not set git credenials helper for %s: %s", remote, err)
				}
			}
		} else {
			glog.Warn("Git directory %s does not have remote origin set", dir)
		}
	}

	traceGitconfig()

	if err, _ := os.ExecCmd(os.ExecParams{Dir: dir, Env: Env}, "git", "pull", "-q"); err != nil {
		return err
	}
	return nil
}
