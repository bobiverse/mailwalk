package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// Execute bash command with valid timeout
// kill process if too long execution
func runBashWithTimeout(timeout time.Duration, cmdstr, cmdInterpreter string) ([]byte, error, error) {
	// arr := strings.Fields(cmdstr)
	// name := arr[0]
	// args := arr[1:]

	// Run command in env as whole
	// Useful when need to execute command with wildcard so these characters
	// is not treated as string
	cmdInterpreter = strings.TrimSpace(cmdInterpreter)
	if cmdInterpreter == "" {
		cmdInterpreter = os.Getenv("SHELL")
		if cmdInterpreter == "" {
			cmdInterpreter = "/bin/bash"
		}
	}
	name := cmdInterpreter
	args := []string{
		"-c",
		strings.TrimPrefix(cmdstr, cmdInterpreter+" -c "), // as one argument
	}

	// color.Magenta("%s %s", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...) // #nosec G204
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	bufOut := &bytes.Buffer{}
	cmd.Stdout = bufOut

	bufErr := &bytes.Buffer{}
	cmd.Stderr = bufErr

	if err := cmd.Start(); err != nil {
		log.Printf("cmd start: %s", err)
	}
	go func() {
		time.Sleep(timeout) // wait in background

		if cmd == nil || cmd.Process == nil {
			return
		}

		pgid, err := syscall.Getpgid(cmd.Process.Pid)
		if err == nil {
			log.Printf("[ KILL ] Kill process of command: %s", name)
			if err := syscall.Kill(-pgid, 15); err != nil { // note the minus sign
				log.Printf("syscall Kill: %s", err)
			}
			// err := cmd.Process.Kill() -- doesnt work on process children
		}

	}()

	err := cmd.Wait()

	var errBuf error
	if _errBuf := bufErr.Bytes(); len(_errBuf) > 0 {
		errBuf = fmt.Errorf("%s", bufErr.Bytes())
	}
	return bufOut.Bytes(), errBuf, err
}
