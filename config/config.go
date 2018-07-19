package config

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
)

// Pane contains a pane configuration
type Pane struct {
	Dir     string   `toml:"dir"`
	Zoom    bool     `toml:"zoom"`
	Split   string   `toml:"split"`
	Scripts []string `toml:"scripts"`
}

// Window contains a window configuration
type Window struct {
	Name        string   `toml:"name"`
	Dir         string   `toml:"dir"`
	Layout      string   `toml:"layout"`
	Sync        bool     `toml:"sync"`
	Scripts     []string `toml:"scripts"`
	PaneScripts []string `toml:"pane-scripts"`
	Panes       []Pane   `toml:"panes"`
}

// Session contains a tmux session configuration
type Session struct {
	Name          string   `toml:"name"`
	Dir           string   `toml:"dir"`
	ClearPanes    bool     `toml:"clear-panes"`
	SelectWindow  string   `toml:"select-window"`
	SelectPane    int      `toml:"select-pane"`
	WindowScripts []string `toml:"window-scripts"`
	Windows       []Window `toml:"windows"`
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
