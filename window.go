package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

type pane struct {
	Dir    string
	Window *window
}

func (p *pane) getDir() string {
	switch {
	case p.Dir != "":
		return p.Dir
	case p.Window.Dir != "":
		return p.Window.Dir
	case p.Window.Sess.Dir != "":
		return p.Window.Sess.Dir
	default:
		return ""
	}
}

type window struct {
	Sess   *session
	Name   string
	Dir    string
	Layout string
	Panes  []*pane
}

func newWindow(sess *session, config windowConfig) *window {
	win := &window{
		Sess:   sess,
		Name:   config.Name,
		Dir:    config.Dir,
		Layout: config.Layout,
	}

	if config.Layout == "" {
		win.Layout = "tiled"
	}

	if config.Dir == "" {
		win.Dir = sess.Dir
	}

	for _, paneConfig := range config.Panes {
		win.Panes = append(win.Panes, &pane{
			Dir:    paneConfig.Dir,
			Window: win,
		})
	}

	return win
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
	if len(w.Panes) <= 1 {
		return nil
	}

	for _, pane := range w.Panes[1:] {
		dir := pane.getDir()
		cmd := exec.Command("tmux", "split-window", "-t", w.Sess.Name+":"+w.Name, "-c", dir)
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
