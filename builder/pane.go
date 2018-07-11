package builder

import (
	"strconv"

	"github.com/alexandrebodin/tmuxctl/config"
	"github.com/alexandrebodin/tmuxctl/tmux"
)

type pane struct {
	Dir     string
	Zoom    bool
	Split   string
	Scripts []string
	Window  *window
	Target  string
}

func newPane(win *window, config config.Pane, index int) *pane {
	normalizedIndex := strconv.Itoa(index + win.Sess.TmuxOptions.PaneBaseIndex)
	pane := &pane{
		Window:  win,
		Zoom:    config.Zoom,
		Split:   config.Split,
		Scripts: config.Scripts,
		Target:  win.Target + "." + normalizedIndex,
	}

	if config.Dir != "" {
		pane.Dir = config.Dir
	} else {
		pane.Dir = win.Dir
	}
	return pane
}

func (p *pane) selectPane() error {
	_, err := tmux.Exec("select-pane", "-t", p.Target)
	return err
}
