package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/alexandrebodin/tmuxctl/tmux"
)

func main() {
	args := []string{".tmuxctlrc"}

	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	filePath := args[0]

	var conf sessionConfig
	if _, err := toml.DecodeFile(filePath, &conf); err != nil {
		log.Fatalf("Error loading configuration %v\n", err)
	}

	runningSessions, err := tmux.ListSessions()
	if err != nil {
		log.Fatal(err)
	}

	options, err := tmux.GetOptions()
	if err != nil {
		log.Fatal(err)
	}

	sess := newSession(conf)
	sess.TmuxOptions = options

	if _, ok := runningSessions[sess.Name]; ok {
		log.Fatalf("Session %s is already running", sess.Name)
	}

	err = sess.start()

	if err != nil {
		log.Fatalf("Error starting session %v\n", err)
	}

	if conf.SelectWindow != "" {
		_, err := tmux.Exec("select-window", "-t", sess.Name+":"+conf.SelectWindow)

		if err != nil {
			log.Fatalf("Error selecting window %s: %v\n", conf.SelectWindow, err)
		}

		if conf.SelectPane != "" {
			_, err := tmux.Exec("select-pane", "-t", sess.Name+":"+conf.SelectWindow+"."+conf.SelectPane)

			if err != nil {
				log.Fatalf("Error selecting pane %s: %v\n", conf.SelectPane, err)
			}
		}
	}

	sess.attach()
}
