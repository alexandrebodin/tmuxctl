package main

type paneConfig struct {
	Dir     string
	Zoom    bool
	Split   string
	Scripts []string
}

type windowConfig struct {
	Name        string
	Dir         string
	Layout      string
	Sync        bool
	Scripts     []string
	Panes       []paneConfig
	PaneScripts []string `toml:"pane-scripts"`
}

type sessionConfig struct {
	Name          string
	Dir           string
	ClearPanes    bool `toml:"clear-panes"`
	Windows       []windowConfig
	SelectWindow  string   `toml:"select-window"`
	SelectPane    int      `toml:"select-pane"`
	WindowScripts []string `toml:"window-scripts"`
}
