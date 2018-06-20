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

	if _, err := os.Stat(conf.Dir); err != nil {
		log.Fatalf("Error with session directory %v\n", err)
	}

	sess := newSession(conf)

	runningSessions, err := tmux.ListSessions()

	if err != nil {
		log.Fatalf("Error getting tmux status %v\n", err)
	}

	if _, ok := runningSessions[sess.Name]; ok {
		log.Fatalf("Session %s is already running", sess.Name)
	}

	err = sess.start()

	if err != nil {
		log.Fatalf("Error starting session %v\n", err)
	}

	sess.attach()
}
