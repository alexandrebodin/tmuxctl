package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/alexandrebodin/tmuxctl/tmux"
)

type session struct {
	Name    string
	Dir     string
	Windows []*window
}

func newSession(config sessionConfig) *session {

	sess := &session{
		Name: config.Name,
		Dir:  lookupDir(config.Dir),
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
	firstWindow := sess.Windows[0]
	_, err := tmux.Exec("new-session", "-d", "-s", sess.Name, "-c", sess.Dir, "-n", firstWindow.Name)
	if err != nil {
		return err
	}

	if firstWindow.Dir != sess.Dir {
		_, err := tmux.Exec("send-keys", "-t", sess.Name+":"+firstWindow.Name, "cd "+firstWindow.Dir, "C-m")
		if err != nil {
			return err
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
		win.renderPane()
		win.renderLayout()
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
