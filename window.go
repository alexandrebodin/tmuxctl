package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

type window struct {
	Sess *session
	Name string
	Dir  string
}

func (w *window) start() error {
	cmd := exec.Command("tmux", "new-window", "-t", w.Sess.Name, "-n", w.Name, "-c", w.Dir)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	runError := cmd.Run()
	if runError != nil {
		return fmt.Errorf("Error Creating tmux session: %v, %q", runError, stderr.String())
	}

	return nil
}
