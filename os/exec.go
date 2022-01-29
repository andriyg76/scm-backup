package os

import (
	"bufio"
	errors2 "errors"
	"fmt"
	log "github.com/andriyg76/glog"
	"io"
	"os"
	"os/exec"
	"time"
)

func read(reader io.Reader, lines *[]string, prefix string, log log.Output) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		log.Printf("%s: %s", prefix, text)
		*lines = append(*lines, text)
	}
}

type ExecParams struct {
	Dir            string
	Ok             []int
	Env            []string
	Stdin          string
	TimoutMinutes  time.Duration
	TimeoutSeconds time.Duration
}

func ExecCmd(params ExecParams, acmd string, args ...string) (error, []string) {
	var timeout = params.TimoutMinutes*time.Minute + params.TimeoutSeconds*time.Second
	if timeout == 0 {
		timeout = 5 * time.Minute
	}

	cmd := exec.Command(acmd, args...)
	cmd.Dir = params.Dir
	log.Trace("command: %s params: %s timeout: %s dir: %s env: %s stdin: %s", acmd, args, timeout, cmd.Dir, cmd.Env, cmd.Stdin)
	if len(params.Env) > 0 {
		cmd.Env = append(os.Environ(), params.Env...)
	}

	log.Debug("executing command %s", acmd)
	stderr, err := cmd.StderrPipe()
	if nil != err {
		log.Fatal("Error obtaining stderr: %s", err.Error())
	}
	stdout, err := cmd.StdoutPipe()
	if nil != err {
		log.Fatal("Error obtaining stdout: %s", err.Error())
	}
	if params.Stdin != "" {
		stdin, err := cmd.StdinPipe()
		if nil != err {
			log.Fatal("Error obtaining stdin: %s", err.Error())
		}
		go func() {
			log.Trace("%s %s", acmd+" in", params.Stdin)
			fmt.Fprintln(stdin, params.Stdin)
			stdin.Close()
		}()
	}

	defer func(pipes ...io.Closer) {
		for _, f := range pipes {
			f.Close()
		}
	}(
		stderr,
		stdout,
	)

	var lines, errors []string
	go read(stdout, &lines, acmd+" out", log.OutputLevel(log.TRACE))
	go read(stderr, &errors, acmd+" err", log.OutputLevel(log.DEBUG))

	if err := cmd.Start(); err != nil {
		err = fmt.Errorf("error starting program: %s, %v", cmd.Path, err.Error())
		log.Error("%s", err)
		return err, nil
	}

	var result = make(chan error)
	go func() {
		result <- cmd.Wait()
		close(result)
	}()

	// Start a timer
	timer := time.After(timeout)
	select {
	case err := <-result:
		log.Trace("Process finished, error: %s", err)
		log.Debug("%s Read text: %v", cmd.Path, lines)
		if nil != err {
			var exiterr *exec.ExitError
			if errors2.As(err, &exiterr) {
				for _, c := range params.Ok {
					if c == exiterr.ExitCode() {
						return nil, lines
					}
				}
			}
			err := log.Error("error starting program: %s, %v", cmd.Path, err.Error())
			log.Error("Stderr: %s", errors)
			return err, lines
		}
		return nil, lines
	case <-timer:
		// Timeout happened first, kill the process and print a message.
		cmd.Process.Kill()
		err := log.Error("command %s timed out, %s", cmd.Path, err)
		log.Error("Stderr: %s", errors)
		log.Info("Read text: %v", lines)
		return err, lines
	}
}
