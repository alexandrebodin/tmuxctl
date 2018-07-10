package config

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
)

// Pane contains a pane configuration
type Pane struct {
	Dir     string
	Zoom    bool
	Split   string
	Scripts []string
}

// Window contains a window configuration
type Window struct {
	Name        string
	Dir         string
	Layout      string
	Sync        bool
	Scripts     []string
	Panes       []Pane
	PaneScripts []string `toml:"pane-scripts"`
}

// Session contains a tmux session configuration
type Session struct {
	Name          string
	Dir           string
	ClearPanes    bool `toml:"clear-panes"`
	Windows       []Window
	SelectWindow  string   `toml:"select-window"`
	SelectPane    int      `toml:"select-pane"`
	WindowScripts []string `toml:"window-scripts"`
}

var (
	validLayouts = []string{
		"even-horizontal",
		"even-vertical",
		"main-horizontal",
		"main-vertical",
		"tiled",
	}
)

func checkValid(conf Session) error {
	// check select-window and select-pane exist
	if conf.SelectWindow != "" {
		var win Window
		found := false
		for _, w := range conf.Windows {
			if w.Name == conf.SelectWindow {
				win = w
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("selected window %s doesn't exist", conf.SelectWindow)
		}

		if conf.SelectPane != 0 {
			if len(win.Panes) < conf.SelectPane {
				return fmt.Errorf("selected pane %d doesn't exist", conf.SelectPane)
			}
		}
	}

	for _, w := range conf.Windows {
		if w.Layout != "" {
			found := false
			for _, l := range validLayouts {
				if l == w.Layout {
					found = true
				}
			}

			if !found {
				return fmt.Errorf("invalid layout '%s' in window '%s'", w.Layout, w.Name)
			}
		}
	}

	// check only on zoom in a window

	return nil
}

// Parse return a sessionConfig from a io.Reader
func Parse(reader io.ReadCloser) (Session, error) {
	defer reader.Close()

	var conf Session
	if _, err := toml.DecodeReader(reader, &conf); err != nil {
		return conf, fmt.Errorf("parsing configuration: %v", err)
	}

	if err := checkValid(conf); err != nil {
		return conf, fmt.Errorf("invalid configuration: %v", err)
	}

	return conf, nil
}
