package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type session struct {
	Name    string
	Dir     string
	Windows []*window
}

func (sess *session) addWindow(w *window) {
	sess.Windows = append(sess.Windows, w)
}

func (sess *session) start() error {
	firstWindow := sess.Windows[0]
	cmd := exec.Command("tmux", "new-session", "-d", "-s", sess.Name, "-c", firstWindow.Dir, "-n", firstWindow.Name)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	runError := cmd.Run()
	if runError != nil {
		return fmt.Errorf("Error Creating tmux session: %v, %q", runError, stderr.String())
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

	cmd = exec.Command("tmux", "select-window", "-t", sess.Name+":"+sess.Windows[0].Name)
	cmd.Stderr = &stderr
	runError = cmd.Run()
	if runError != nil {
		return fmt.Errorf("Error Creating tmux session: %v, %q", runError, stderr.String())
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
