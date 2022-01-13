package main

import (
	"github.com/andiryg76/scm_backup/os"
	log "github.com/andriyg76/glogger"
)

func main() {
	log.SetLevel(log.DEBUG)

	os.ExecCmd("bash", "-c", "false")
}
