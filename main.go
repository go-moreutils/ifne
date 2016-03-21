package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "Usage: ifne [-n] command [args]")
		os.Exit(255)
	}

	do := true
	if os.Args[1] == "-n" {
		do = false
		os.Args = os.Args[:1+copy(os.Args[1:], os.Args[2:])]
	}

	b, _ := ioutil.ReadAll(os.Stdin)
	if (do == true && len(b) == 0) || (do == false && len(b) > 0) {
		return
	}
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Stdin = bytes.NewBuffer(b)
	b, err := cmd.CombinedOutput()
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "%s: %s\n", err, string(b))
	status := 1
	state := cmd.ProcessState
	if state != nil {
		if ws, ok := state.Sys().(syscall.WaitStatus); ok {
			if ws.Exited() {
				status = ws.ExitStatus()
			}
		}
	}
	os.Exit(status)
}
