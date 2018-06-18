package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type windowConfig struct {
	Dir string
}

type config struct {
	Name    string
	Dir     string
	Windows map[string]*windowConfig
}

func main() {
	args := []string{".tmuxrc"}

	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	filePath := args[0]

	var conf config
	if _, err := toml.DecodeFile(filePath, &conf); err != nil {
		panic(fmt.Errorf("Error decoding configuration %s", err))
	}

	if _, err := os.Stat(conf.Dir); err != nil {
		log.Fatal(err)
	}

	sess := &session{
		Name: conf.Name,
		Dir:  conf.Dir,
	}

	for winName, winConfig := range conf.Windows {
		window := &window{
			Sess: sess,
			Name: winName,
			Dir:  winConfig.Dir,
		}
		sess.addWindow(window)
	}

	err := sess.start()

	if err != nil {
		log.Fatal(err)
	}

	sess.attach()
}
