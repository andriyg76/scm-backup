package os

import (
	"github.com/andriyg76/glogger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecTimout(t *testing.T) {
	glogger.SetLevel(glogger.TRACE)

	err := execCmdInt(params{timeoutSeconds: 1}, "ping", "localhost")
	assert.Error(t, err)
}

func TestExecNoTimout(t *testing.T) {
	glogger.SetLevel(glogger.TRACE)

	err := execCmdInt(params{timeoutSeconds: 3}, "bash", "-x", "-c", "sleep 1")
	assert.Nil(t, err)
}

func TestExecFailure(t *testing.T) {
	glogger.SetLevel(glogger.TRACE)

	err := execCmdInt(params{timeoutSeconds: 3}, "bash", "-x", "-c", "false")
	assert.Error(t, err)
}

func TestExecOkNoNull(t *testing.T) {
	glogger.SetLevel(glogger.TRACE)

	err := execCmdInt(params{timeoutSeconds: 3, ok: []int{1}}, "bash", "-x", "-c", "false")
	assert.Error(t, err)
}
