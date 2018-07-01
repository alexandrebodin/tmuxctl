package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/alexandrebodin/go-findup"
)

func main() {
	var filePath string
	var err error
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	if filePath == "" {
		filePath, err = findup.Find(".tmuxctlrc")
		if err != nil {
			log.Fatal(err)
		}
	}

	if strings.Trim(filePath, " ") == "" {
		log.Fatal("not file path provided")
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
