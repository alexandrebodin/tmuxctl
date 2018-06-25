package main

type paneConfig struct {
	Dir   string
	Zoom  bool
	Split string
}

type windowConfig struct {
	Name   string
	Dir    string
	Layout string
	Sync   bool
	Panes  []paneConfig
}

type sessionConfig struct {
	Name         string
	Dir          string
	Windows      []windowConfig
	SelectWindow string `toml:"select-window"`
	SelectPane   int    `toml:"select-pane"`
}
