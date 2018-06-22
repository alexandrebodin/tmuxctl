package main

type paneConfig struct {
	Dir string
}

type windowConfig struct {
	Name   string
	Dir    string
	Layout string
	Panes  []paneConfig
}

type sessionConfig struct {
	Name         string
	Dir          string
	Windows      []windowConfig
	SelectWindow string `toml:"select-window"`
	SelectPane   string `toml:"select-pane"`
}
