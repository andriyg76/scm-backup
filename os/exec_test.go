package os

import (
	"github.com/andriyg76/glogger"
	list2 "github.com/andriyg76/scm-backup/lists"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecTimout(t *testing.T) {
	glogger.SetLevel(glogger.TRACE)

	err, _ := execCmdInt(intParams{timeoutSeconds: 1}, "ping", "localhost")
	assert.Error(t, err)
}

func TestExecNoTimout(t *testing.T) {
	glogger.SetLevel(glogger.TRACE)

	err, _ := execCmdInt(intParams{timeoutSeconds: 3}, "bash", "-x", "-c", "sleep 1")
	assert.Nil(t, err)
}

func TestExecFailure(t *testing.T) {
	glogger.SetLevel(glogger.TRACE)

	err, _ := execCmdInt(intParams{timeoutSeconds: 3}, "bash", "-x", "-c", "false")
	assert.Error(t, err)
}

func TestExecOkNoNull(t *testing.T) {
	glogger.SetLevel(glogger.TRACE)

	err, _ := execCmdInt(intParams{timeoutSeconds: 3, ExecParams: ExecParams{Ok: []int{1}}}, "bash", "-x", "-c", "false")
	assert.Nil(t, err)
}

func TestStdin(t *testing.T) {
	err, list := ExecCmd(ExecParams{Stdin: "ping-pong"}, "cat")
	assert.Nil(t, err)

	assert.Equal(t, list2.String("ping-pong"), list)
}
func TestExecData(t *testing.T) {
	glogger.SetLevel(glogger.TRACE)

	err, list := execCmdInt(intParams{}, "bash", "-x", "-c",
		"echo 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 - 1  -;"+
			"echo 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 - 2  -;"+
			"sleep 0;"+
			"echo 1234567890 - end -",
	)
	assert.Nil(t, err)
	assert.Equal(t, list2.String("1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 - 1 -",
		"1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 - 2 -",
		"1234567890 - end -"), list)

	err, list = ExecCmd(ExecParams{}, "bash", "-x", "-c", "echo 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890; echo 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890; echo 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890; echo 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 echo 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890; echo 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890; echo 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 echo 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890; echo 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890; echo 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890; sleep 1; echo -e the end -")
	assert.Nil(t, err)
	assert.Equal(t, list2.String("1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890",
		"1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890",
		"1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890",
		"1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 echo 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890",
		"1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890",
		"1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 echo 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890",
		"1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890",
		"1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890",
		"the end -"), list)
}
