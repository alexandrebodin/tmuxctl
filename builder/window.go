package builder

import (
	"fmt"
	"strings"

	"github.com/alexandrebodin/tmuxctl/config"
	"github.com/alexandrebodin/tmuxctl/tmux"
)

type window struct {
	Sess        *Session
	Name        string
	Dir         string
	Layout      string
	Sync        bool
	Scripts     []string
	Panes       []*pane
	PaneScripts []string
	Target      string
}

func newWindow(sess *Session, config config.Window) *window {
	win := &window{
		Sess:        sess,
		Name:        config.Name,
		Layout:      config.Layout,
		Sync:        config.Sync,
		Scripts:     config.Scripts,
		PaneScripts: config.PaneScripts,
		Target:      sess.Name + ":" + config.Name,
	}

	if config.Dir != "" {
		win.Dir = lookupDir(config.Dir)
	} else {
		win.Dir = sess.Dir
	}

	for idx, paneConfig := range config.Panes {
		win.Panes = append(win.Panes, newPane(win, paneConfig, idx))
	}

	return win
}

func (w *window) start() error {
	_, err := tmux.Exec("new-window", "-t", w.Sess.Name, "-n", w.Name, "-c", w.Dir)
	if err != nil {
		return fmt.Errorf("starting window: %v", err)
	}
	return nil
}

func (w *window) runScripts() error {
	for _, script := range w.Scripts {
		err := tmux.SendKeys(w.Sess.Name+":"+w.Name, script)
		if err != nil {
			return fmt.Errorf("run window scripts: %v", err)
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
		_, err := tmux.Exec("set-window-option", "-t", w.Target, "synchronize-panes")
		return err
	}

	return nil
}

func (w *window) runPaneScripts() error {
	for _, pane := range w.Panes {
		for _, script := range w.PaneScripts {
			err := tmux.SendKeys(pane.Target, script)
			if err != nil {
				return err
			}
		}

		for _, script := range pane.Scripts {
			err := tmux.SendKeys(pane.Target, script)
			if err != nil {
				return err
			}
		}

		// clearing panes
		if w.Sess.ClearPanes {
			err := tmux.SendRawKeys(pane.Target, "C-l")
			if err != nil {
				return err
			}
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
		err := tmux.SendKeys(firstPane.Target, "cd "+firstPane.Dir)
		if err != nil {
			return err
		}
	}

	for _, pane := range w.Panes[1:] {
		args := []string{"split-window", "-t", w.Target, "-c", pane.Dir}

		if pane.Split != "" {
			args = append(args, strings.Split(pane.Split, " ")...)
		}
		_, err := tmux.Exec(args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *window) renderLayout() error {
	if w.Layout != "" {
		_, err := tmux.Exec("select-layout", "-t", w.Target, w.Layout)
		return err
	}

	return nil
}

func (w *window) zoomPanes() error {
	for _, pane := range w.Panes {
		if pane.Zoom {
			_, err := tmux.Exec("resize-pane", "-t", pane.Target, "-Z")
			if err != nil {
				return err
			}

			// stop after first pane zoomed
			return nil
		}
	}

	return nil
}

func (w *window) selectWindow() error {
	_, err := tmux.Exec("select-window", "-t", w.Target)
	return err
}

func (w *window) selectPane(index int) (*pane, error) {
	if index > len(w.Panes) {
		return nil, fmt.Errorf("pane %d not found", index)
	}

	p := w.Panes[index-1]
	return p, p.selectPane()
}
