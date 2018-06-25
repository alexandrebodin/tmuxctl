package main

type pane struct {
	Dir    string
	Zoom   bool
	Window *window
}

func newPane(win *window, config paneConfig) *pane {
	pane := &pane{
		Window: win,
		Zoom:   config.Zoom,
	}

	if config.Dir != "" {
		pane.Dir = lookupDir(config.Dir)
	} else {
		pane.Dir = win.Dir
	}
	return pane
}
