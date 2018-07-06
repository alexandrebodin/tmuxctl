package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/alexandrebodin/tmuxctl/config"
	"github.com/alexandrebodin/tmuxctl/tmux"
)

type session struct {
	TmuxOptions   *tmux.Options
	Name          string
	Dir           string
	Windows       []*window
	ClearPanes    bool
	WindowScripts []string
}

func newSession(config config.Session, options *tmux.Options) *session {
	sess := &session{
		Name:          config.Name,
		Dir:           lookupDir(config.Dir),
		ClearPanes:    config.ClearPanes,
		WindowScripts: config.WindowScripts,
		TmuxOptions:   options,
	}

	for _, winConfig := range config.Windows {
		window := newWindow(sess, winConfig)
		sess.Windows = append(sess.Windows, window)
	}

	return sess
}

func (sess *session) start() error {
	// get term size
	width, height, err := getTermSize()
	if err != nil {
		return err
	}

	if len(sess.Windows) == 0 {
		return errors.New("session has no window")
	}

	firstWindow := sess.Windows[0]
	_, err = tmux.Exec("new-session", "-d", "-s", sess.Name, "-c", sess.Dir, "-n", firstWindow.Name, "-x", width, "-y", height)
	if err != nil {
		return fmt.Errorf("starting session: %v", err)
	}

	if firstWindow.Dir != sess.Dir {
		cdCmd := fmt.Sprintf("cd %s", firstWindow.Dir)
		err := tmux.SendKeys(sess.Name+":"+firstWindow.Name, cdCmd)
		if err != nil {
			return fmt.Errorf("moving window to dir %s: %v", firstWindow.Dir, err)
		}
	}

	if len(sess.Windows) > 1 {
		for _, win := range sess.Windows[1:] {
			err := win.start()
			if err != nil {
				return fmt.Errorf("starting window %s: %v", win.Name, err)
			}
		}
	}

	for _, win := range sess.Windows {

		for _, script := range sess.WindowScripts {
			err := tmux.SendKeys(sess.Name+":"+win.Name, script)
			if err != nil {
				return err
			}
		}

		err := win.init()
		if err != nil {
			return fmt.Errorf("initializing window: %v", err)
		}
	}

	return nil
}

func (sess *session) attach() error {
	tmux, err := exec.LookPath("tmux")
	if err != nil {
		return fmt.Errorf("looking up tmux: %v", err)
	}

	args := []string{"tmux", "attach", "-t", sess.Name}
	if sysErr := syscall.Exec(tmux, args, os.Environ()); sysErr != nil {
		return fmt.Errorf("attaching to session: %v", sysErr)
	}

	return nil
}

func (sess *session) selectWindow(name string) (*window, error) {
	for _, w := range sess.Windows {
		if w.Name == name {
			return w, w.selectWindow()
		}
	}

	return nil, fmt.Errorf("window %s not found", name)
}
