package main

import (
	"strconv"

	"github.com/alexandrebodin/tmuxctl/tmux"
)

type pane struct {
	Dir    string
	Window *window
}

func newPane(win *window, config paneConfig) *pane {
	pane := &pane{
		Window: win,
	}

	if config.Dir != "" {
		pane.Dir = lookupDir(config.Dir)
	} else {
		pane.Dir = win.Dir
	}
	return pane
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
		Layout: config.Layout,
	}

	if config.Dir != "" {
		win.Dir = lookupDir(config.Dir)
	} else {
		win.Dir = sess.Dir
	}

	if config.Layout == "" {
		win.Layout = "tiled"
	}

	for _, paneConfig := range config.Panes {
		win.Panes = append(win.Panes, newPane(win, paneConfig))
	}

	return win
}

func (w *window) start() error {
	args := []string{"new-window", "-t", w.Sess.Name, "-n", w.Name, "-c", w.Dir}
	_, err := tmux.Exec(args...)
	return err
}

func (w *window) renderPane() error {
	if len(w.Panes) == 0 {
		return nil
	}

	firstPane := w.Panes[0]
	if firstPane.Dir != "" && firstPane.Dir != w.Dir { // we need to move the pane
		_, err := tmux.Exec("send-keys", "-t", w.Sess.Name+":"+w.Name+"."+strconv.Itoa(w.Sess.TmuxOptions.PaneBaseIndex), "cd "+firstPane.Dir, "C-m")
		if err != nil {
			return err
		}
	}

	for _, pane := range w.Panes[1:] {
		args := []string{"split-window", "-t", w.Sess.Name + ":" + w.Name, "-c", pane.Dir}

		_, err := tmux.Exec(args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *window) renderLayout() error {
	_, err := tmux.Exec("select-layout", "-t", w.Sess.Name+":"+w.Name, w.Layout)
	return err
}
