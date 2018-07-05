package main

import (
	"fmt"
	"log"
	"os"

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
	if err != nil {
		log.Fatal(err)
	}

	options, err := GetOptions()
	if err != nil {
		log.Fatal(err)
	}

	sess := newSession(conf, options)
	if _, ok := runningSessions[sess.Name]; ok {
		log.Fatalf("Session %s is already running", sess.Name)
	}

	checkError := func(err error) {
		if err != nil {
			log.Println(err)
			// kill session if an error occurs after starting it
			Exec("kill-session", "-t", sess.Name)
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
			// index := strconv.Itoa(conf.SelectPane + (options.PaneBaseIndex - 1))
			// target := fmt.Sprintf("%s:%s.%s", sess.Name, conf.SelectWindow, index)
			// _, err := Exec("select-pane", "-t", target)
			checkError(err)
		}
	}

	err = sess.attach()
	checkError(err)
}
