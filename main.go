package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/alexandrebodin/go-findup"
	"github.com/alexandrebodin/tmuxctl/config"
	"github.com/alexandrebodin/tmuxctl/tmux"
)

var (
	version    = "development"
	start      = kingpin.Command("start", "start a tmux instance").Default()
	configPath = start.Arg("config", "Tmux config file").Default(".tmuxctlrc").String()
)

func main() {
	kingpin.Version(fmt.Sprintf("tmuxct %s", version)).Author("Alexandre BODIN")
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.Parse()

	fmt.Printf("Start tmux with config file: %v\n", *configPath)

	filePath, err := findup.Find(*configPath)
	if err != nil {
		log.Fatalf("Error locating config file %v\n", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error openning config file %v\n", err)
	}

	conf, err := config.Parse(file)
	if err != nil {
		log.Fatalf("Error parsing %s: %v\n", filePath, err)
	}

	runningSessions, err := tmux.ListSessions()
	if err != nil {
		log.Fatalf("Error listing running sessions %v\n", err)
	}

	options, err := tmux.GetOptions()
	if err != nil {
		log.Fatalf("Error getting tmux options %v\n", err)
	}

	sess := newSession(conf, options)
	if _, ok := runningSessions[sess.Name]; ok {
		log.Fatalf("Session %s is already running\n", sess.Name)
	}

	checkError := func(err error) {
		if err != nil {
			log.Println(err)
			// kill session if an error occurs after starting it
			tmux.Exec("kill-session", "-t", sess.Name)
			os.Exit(-1)
		}
	}

	err = sess.start()
	checkError(err)

	if conf.SelectWindow != "" {
		w, err := sess.selectWindow(conf.SelectWindow)
		checkError(err)

		if conf.SelectPane != 0 {
			_, err := w.selectPane(conf.SelectPane)
			checkError(err)
		}
	}

	err = sess.attach()
	checkError(err)
}
