package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type session struct {
	TmuxOptions   *Options
	Name          string
	Dir           string
	Windows       []*window
	ClearPanes    bool
	WindowScripts []string
}

func newSession(config sessionConfig) *session {
	sess := &session{
		Name:          config.Name,
		Dir:           lookupDir(config.Dir),
		ClearPanes:    config.ClearPanes,
		WindowScripts: config.WindowScripts,
	}

	for _, winConfig := range config.Windows {
		window := newWindow(sess, winConfig)
		sess.addWindow(window)
	}

	return sess
}

func (sess *session) addWindow(w *window) {
	sess.Windows = append(sess.Windows, w)
}

func (sess *session) start() error {
	// get term size
	width, height, err := getTermSize()
	if err != nil {
		return err
	}

	firstWindow := sess.Windows[0]
	_, err = Exec("new-session", "-d", "-s", sess.Name, "-c", sess.Dir, "-n", firstWindow.Name, "-x", width, "-y", height)
	if err != nil {
		return fmt.Errorf("error starting session %v", err)
	}

	if firstWindow.Dir != sess.Dir {
		cdCmd := fmt.Sprintf("cd %s", firstWindow.Dir)
		err := SendKeys(sess.Name+":"+firstWindow.Name, cdCmd)
		if err != nil {
			return fmt.Errorf("error moving to dir %s, %v", firstWindow.Dir, err)
		}
	}

	if len(sess.Windows) > 1 {
		for _, win := range sess.Windows[1:] {
			err := win.start()
			if err != nil {
				return fmt.Errorf("Error starting window %v", err)
			}
		}
	}

	for _, win := range sess.Windows {

		for _, script := range sess.WindowScripts {
			err := SendKeys(sess.Name+":"+win.Name, script)
			if err != nil {
				return err
			}
		}

		err := win.init()
		if err != nil {
			return fmt.Errorf("Error initializing window %v", err)
		}
	}

	return nil
}

func (sess *session) attach() error {
	tmux, err := exec.LookPath("tmux")
	if err != nil {
		return fmt.Errorf("Error looking up tmux %v", err)
	}

	args := []string{"tmux", "attach", "-t", sess.Name}
	if sysErr := syscall.Exec(tmux, args, os.Environ()); sysErr != nil {
		return fmt.Errorf("Error attaching to session %s: %v", sess.Name, sysErr)
	}

	return nil
}
