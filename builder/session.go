package builder

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/alexandrebodin/tmuxctl/config"
	"github.com/alexandrebodin/tmuxctl/term"
	"github.com/alexandrebodin/tmuxctl/tmux"
)

// Session struct represents a tmux session
type Session struct {
	TmuxOptions   *tmux.Options
	Name          string
	Dir           string
	Windows       []*window
	ClearPanes    bool
	WindowScripts []string
	SelectWindow  string
	SelectPane    int
}

// NewSession create a tmux session instance
func NewSession(config config.Session, options *tmux.Options) *Session {
	sess := &Session{
		Name:          config.Name,
		Dir:           config.Dir,
		ClearPanes:    config.ClearPanes,
		WindowScripts: config.WindowScripts,
		SelectWindow:  config.SelectWindow,
		SelectPane:    config.SelectPane,
		TmuxOptions:   options,
	}

	for idx, winConfig := range config.Windows {
		window := newWindow(sess, winConfig, idx)
		sess.Windows = append(sess.Windows, window)
	}

	return sess
}

// Start starts a tmux sessions
func (sess *Session) Start() error {
	// get term size
	width, height, err := term.GetDimensions()
	if err != nil {
		return err
	}

	args := []string{"new-session", "-d", "-s", sess.Name, "-x", width, "-y", height}
	if len(sess.Windows) == 0 {
		_, err = tmux.Exec()
	} else {
		firstWindow := sess.Windows[0]
		args = append(args, "-n", firstWindow.Name)
	}

	_, err = tmux.Exec(args...)
	if err != nil {
		return fmt.Errorf("starting session: %v", err)
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
		err := tmux.SendKeys(win.Target, fmt.Sprintf("cd %s", win.Dir))
		if err != nil {
			return fmt.Errorf("moving window to dir %s: %v", win.Dir, err)
		}

		for _, script := range sess.WindowScripts {
			err := tmux.SendKeys(win.Target, script)
			if err != nil {
				return fmt.Errorf("running window scripts: %v", err)
			}
		}

		err = win.init()
		if err != nil {
			return fmt.Errorf("initializing window: %v", err)
		}
	}

	if sess.SelectWindow != "" {
		w, err := sess.selectWindow(sess.SelectWindow)
		if err != nil {
			return err
		}

		if sess.SelectPane != 0 {
			_, err := w.selectPane(sess.SelectPane)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Attach attaches the process to the tmux session
func (sess *Session) Attach() error {
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

func (sess *Session) selectWindow(name string) (*window, error) {
	for _, w := range sess.Windows {
		if w.Name == name {
			return w, w.selectWindow()
		}
	}

	return nil, fmt.Errorf("window %s not found", name)
}
