package main

import (
	"strconv"
	"strings"
)

type window struct {
	Sess        *session
	Name        string
	Dir         string
	Layout      string
	Sync        bool
	Scripts     []string
	Panes       []*pane
	PaneScripts []string
}

func newWindow(sess *session, config windowConfig) *window {
	win := &window{
		Sess:        sess,
		Name:        config.Name,
		Layout:      config.Layout,
		Sync:        config.Sync,
		Scripts:     config.Scripts,
		PaneScripts: config.PaneScripts,
	}

	if config.Dir != "" {
		win.Dir = lookupDir(config.Dir)
	} else {
		win.Dir = sess.Dir
	}

	for _, paneConfig := range config.Panes {
		win.Panes = append(win.Panes, newPane(win, paneConfig))
	}

	return win
}

func (w *window) start() error {
	args := []string{"new-window", "-t", w.Sess.Name, "-n", w.Name, "-c", w.Dir}
	_, err := Exec(args...)
	return err
}

func (w *window) runScripts() error {
	for _, script := range w.Scripts {
		err := SendKeys(w.Sess.Name+":"+w.Name, script)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *window) init() error {
	var err error
	err = w.runScripts()
	if err != nil {
		return err
	}

	err = w.renderPane()
	if err != nil {
		return err
	}

	err = w.renderLayout()
	if err != nil {
		return err
	}

	err = w.zoomPanes()
	if err != nil {
		return err
	}

	err = w.runPaneScripts()
	if err != nil {
		return err
	}

	if w.Sync {
		_, err := Exec("set-window-option", "-t", w.Sess.Name+":"+w.Name, "synchronize-panes")
		return err
	}

	return nil
}

func (w *window) runPaneScripts() error {
	for idx, pane := range w.Panes {
		paneTarget := w.Sess.Name + ":" + w.Name + "." + strconv.Itoa(idx+w.Sess.TmuxOptions.PaneBaseIndex)

		for _, script := range w.PaneScripts {
			err := SendKeys(paneTarget, script)
			if err != nil {
				return err
			}
		}

		for _, script := range pane.Scripts {
			err := SendKeys(paneTarget, script)
			if err != nil {
				return err
			}
		}

		// clearing panes
		err := SendRawKeys(paneTarget, "C-l")
		if err != nil {
			return err
		}

	}

	return nil
}

func (w *window) renderPane() error {
	if len(w.Panes) == 0 {
		return nil
	}

	firstPane := w.Panes[0]
	if firstPane.Dir != "" && firstPane.Dir != w.Dir { // we need to move the pane
		err := SendKeys(w.Sess.Name+":"+w.Name+"."+strconv.Itoa(w.Sess.TmuxOptions.PaneBaseIndex), "cd "+firstPane.Dir)
		if err != nil {
			return err
		}
	}

	for _, pane := range w.Panes[1:] {
		args := []string{"split-window", "-t", w.Sess.Name + ":" + w.Name, "-c", pane.Dir}

		if pane.Split != "" {
			args = append(args, strings.Split(pane.Split, " ")...)
		}
		_, err := Exec(args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *window) renderLayout() error {
	if w.Layout != "" {
		_, err := Exec("select-layout", "-t", w.Sess.Name+":"+w.Name, w.Layout)
		return err
	}

	return nil
}

func (w *window) zoomPanes() error {
	for idx, pane := range w.Panes {
		if pane.Zoom {
			index := strconv.Itoa(idx + w.Sess.TmuxOptions.PaneBaseIndex)
			_, err := Exec("resize-pane", "-t", w.Sess.Name+":"+w.Name+"."+index, "-Z")
			if err != nil {
				return err
			}

			return nil // stop after first pane zoomed
		}
	}

	return nil
}
