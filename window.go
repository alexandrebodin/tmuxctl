package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

type pane struct {
}

type window struct {
	Sess   *session
	Name   string
	Dir    string
	Layout string
	Panes  []*pane
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

func (w *window) renderPane() error {
	for range w.Panes {
		cmd := exec.Command("tmux", "split-window", "-t", w.Sess.Name+":"+w.Name)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		runError := cmd.Run()
		if runError != nil {
			return fmt.Errorf("Error Creating tmux session: %v, %q", runError, stderr.String())
		}
	}

	return nil
}

func (w *window) renderLayout() error {
	cmd := exec.Command("tmux", "select-layout", "-t", w.Sess.Name+":"+w.Name, w.Layout)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	runError := cmd.Run()
	if runError != nil {
		return fmt.Errorf("Error Creating tmux session: %v, %q", runError, stderr.String())
	}
	return nil
}
