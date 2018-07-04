package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/alecthomas/kingpin"
	"github.com/alexandrebodin/go-findup"
)

var (
	version = "development"
	start   = kingpin.Command("start", "start a tmux instance").Default()
	config  = start.Arg("config", "Tmux config file").Default(".tmuxctlrc").String()
)

func main() {
	kingpin.Version(fmt.Sprintf("tmuxct %s", version)).Author("Alexandre BODIN")
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.Parse()

	fmt.Printf("Start tmux with config file: %v\n", *config)

	filePath, err := findup.Find(*config)
	if err != nil {
		log.Fatal(err)
	}

	var conf sessionConfig
	if _, err := toml.DecodeFile(filePath, &conf); err != nil {
		log.Fatalf("Error loading configuration %v\n", err)
	}

	runningSessions, err := ListSessions()
	checkError(err)

	options, err := GetOptions()
	checkError(err)

	sess := newSession(conf)
	sess.TmuxOptions = options

	if _, ok := runningSessions[sess.Name]; ok {
		log.Fatalf("Session %s is already running", sess.Name)
	}

	err = sess.start()
	checkError(err)

	if conf.SelectWindow != "" {
		_, err := Exec("select-window", "-t", sess.Name+":"+conf.SelectWindow)
		checkError(err)

		if conf.SelectPane != 0 {
			index := strconv.Itoa(conf.SelectPane + (options.PaneBaseIndex - 1))
			_, err := Exec("select-pane", "-t", sess.Name+":"+conf.SelectWindow+"."+index)
			checkError(err)
		}
	}

	err = sess.attach()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
