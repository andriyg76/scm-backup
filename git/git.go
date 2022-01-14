package git

import (
	"github.com/andiryg76/scm_backup/lists"
	"github.com/andiryg76/scm_backup/os"
	"github.com/andriyg76/glogger"
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
	return nil
}
