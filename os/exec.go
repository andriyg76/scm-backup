package os

import (
	"bufio"
	"fmt"
	log "github.com/andriyg76/glogger"
	"io"
	"os/exec"
	"time"
)

func read(reader io.Reader, lines *[]string, log func(format string, objs ...interface{})) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		log("git: %s", text)
		*lines = append(*lines, text)
	}
}

type params struct {
	timoutMinutes  time.Duration
	timeoutSeconds time.Duration
	ok             []int
}

func ExecCmd(acmd string, args ...string) error {
	return execCmdInt(params{
		timoutMinutes: 5,
	}, acmd, args...)
}

func execCmdInt(params params, acmd string, args ...string) error {
	var timeout = params.timoutMinutes*time.Minute + params.timeoutSeconds*time.Second
	if timeout == 0 {
		timeout = 5 * time.Minute
	}
	log.Trace("Timeout: %s", timeout)

	cmd := exec.Command(acmd, args...)
	log.Debug("executing command %s", acmd)
	stderr, err := cmd.StderrPipe()
	if nil != err {
		log.Fatal("Error obtaining stdin: %s", err.Error())
	}
	stdout, err := cmd.StdoutPipe()
	if nil != err {
		log.Fatal("Error obtaining stdout: %s", err.Error())
	}

	defer func(pipes ...io.ReadCloser) {
		for _, f := range pipes {
			f.Close()
		}
	}(
		stderr,
		stdout,
	)

	var lines, errors []string
	go read(stdout, &lines, log.Debug)
	go read(stderr, &errors, log.Error)

	if err := cmd.Start(); err != nil {
		err = fmt.Errorf("error starting program: %s, %v", cmd.Path, err.Error())
		log.Error("%s", err)
		return err
	}

	var result chan error
	go func() {
		result <- cmd.Wait()
	}()

	// Start a timer
	timer := time.After(timeout)
	select {
	case err := <-result:
		log.Info("Read text: %s", lines)
		if nil != err {
			err := fmt.Errorf("Error starting program: %s, %v", cmd.Path, err.Error())
			if exiterr, ok := err.(*exec.ExitError); ok {
				for _, c := range params.ok {
					if c == exiterr.ExitCode() {
						return nil
					}
				}
			}
			log.Error("%s", err)
			return err
		}
		return nil
	case <-timer:
		// Timeout happened first, kill the process and print a message.
		cmd.Process.Kill()
		return fmt.Errorf("command timed out")
	}
}
