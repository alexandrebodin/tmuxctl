package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/alexandrebodin/go-findup"
	"github.com/alexandrebodin/tmuxctl/builder"
	"github.com/alexandrebodin/tmuxctl/config"
	"github.com/alexandrebodin/tmuxctl/tmux"
)

func init() {
	var start = rootCmd.Command("start", "Start a tmux instance").Default()
	configPath := start.Arg("config", "Tmux config file").Default(".tmuxctlrc").String()
	test := start.Flag("test", "Test config file").Short('t').Bool()

	start.Action(func(ctx *kingpin.ParseContext) error {

		filePath, err := findup.Find(*configPath)
		if err != nil {
			kingpin.FatalUsageContext(ctx, "locating config file: %v", err)
		}

		file, err := os.Open(filePath)
		if err != nil {
			kingpin.FatalUsageContext(ctx, "Error openning config file %v\n", err)
		}

		conf, err := config.Parse(file)
		if err != nil {
			kingpin.FatalUsageContext(ctx, "Error parsing %s: %v\n", file.Name(), err)
		}

		// stop after parsing in case of test
		if *test {
			fmt.Printf("Config file %s is valid\n", file.Name())
			os.Exit(0)
		}

		fmt.Printf("Start tmux with config file: %v\n", *configPath)

		runningSessions, err := tmux.ListSessions()
		if err != nil {
			kingpin.FatalUsageContext(ctx, "Error listing running sessions %v\n", err)
		}

		options, err := tmux.GetOptions()
		if err != nil {
			kingpin.FatalUsageContext(ctx, "Error getting tmux options %v\n", err)
		}

		sess := builder.NewSession(conf, options)
		if _, ok := runningSessions[sess.Name]; ok {
			kingpin.FatalUsageContext(ctx, "Session %s is already running\n", sess.Name)
		}

		checkError := func(err error) {
			if err != nil {
				log.Println(err)
				// kill session if an error occurs after starting it
				tmux.Exec("kill-session", "-t", sess.Name)
				os.Exit(-1)
			}
		}

		err = sess.Start()
		checkError(err)

		err = sess.Attach()
		checkError(err)
		return nil
	})
}
