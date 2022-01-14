package os

import (
	"bufio"
	errors2 "errors"
	"fmt"
	log "github.com/andriyg76/glogger"
	"io"
	"os/exec"
	"strings"
	"time"
)

func read(reader io.Reader, lines *[]string, prefix string, log func(format string, objs ...interface{})) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		log("%s: %s", prefix, text)
		*lines = append(*lines, text)
	}
}

type ExecParams struct {
	Dir   string
	Ok    []int
	Env   []string
	Stdin string
}

type intParams struct {
	timoutMinutes  time.Duration
	timeoutSeconds time.Duration
	ExecParams
}

func ExecCmd(params ExecParams, acmd string, args ...string) (error, []string) {
	return execCmdInt(intParams{
		timoutMinutes: 5,
		ExecParams:    params,
	}, acmd, args...)
}

func execCmdInt(params intParams, acmd string, args ...string) (error, []string) {
	var timeout = params.timoutMinutes*time.Minute + params.timeoutSeconds*time.Second
	if timeout == 0 {
		timeout = 5 * time.Minute
	}

	cmd := exec.Command(acmd, args...)
	cmd.Dir = params.Dir
	cmd.Env = params.Env

	log.Trace("command: %s params: %s timeout: %s dir: %s env: %s stdin: %s", acmd, args, timeout, cmd.Dir, cmd.Env, cmd.Stdin)

	log.Debug("executing command %s", acmd)
	stderr, err := cmd.StderrPipe()
	if nil != err {
		log.Fatal("Error obtaining stderr: %s", err.Error())
	}
	stdout, err := cmd.StdoutPipe()
	if nil != err {
		log.Fatal("Error obtaining stdout: %s", err.Error())
	}
	cmd.Stdin = strings.NewReader(params.Stdin)

	defer func(pipes ...io.Closer) {
		for _, f := range pipes {
			f.Close()
		}
	}(
		stderr,
		stdout,
	)

	var lines, errors []string
	go read(stdout, &lines, acmd, log.Debug)
	go read(stderr, &errors, acmd, log.Error)

	if err := cmd.Start(); err != nil {
		err = fmt.Errorf("error starting program: %s, %v", cmd.Path, err.Error())
		log.Error("%s", err)
		return err, nil
	}

	var result chan error = make(chan error)
	go func() {
		result <- cmd.Wait()
		close(result)
	}()

	// Start a timer
	timer := time.After(timeout)
	select {
	case err := <-result:
		log.Trace("Process finished, error: %s", err)
		log.Info("Read text: %v", lines)
		if nil != err {
			var exiterr *exec.ExitError
			if errors2.As(err, &exiterr) {
				for _, c := range params.Ok {
					if c == exiterr.ExitCode() {
						return nil, lines
					}
				}
			}
			err := fmt.Errorf("error starting program: %s, %v", cmd.Path, err.Error())
			log.Error("%s", err)
			return err, lines
		}
		return nil, lines
	case <-timer:
		// Timeout happened first, kill the process and print a message.
		cmd.Process.Kill()
		err := fmt.Errorf("command %s timed out", cmd.Path)
		log.Error("Timeout happened %s", err)
		log.Info("Read text: %v", lines)
		return err, lines
	}
}
