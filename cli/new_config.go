package cli

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/alecthomas/kingpin"
	"github.com/alexandrebodin/tmuxctl/config"
)

func init() {
	var newConfig = rootCmd.Command("new", "Create a new configuration")
	filePtr := newConfig.Flag("file", "Path of the configuration file").Short('f').Default(".tmuxctlrc").String()
	namePtr := newConfig.Flag("name", "Name of the session").Short('n').String()
	wCount := newConfig.Flag("windows", "Count of windows").Short('w').Default("1").Int()
	pCount := newConfig.Flag("panes", "Count of panes per window").Short('p').Default("1").Int()

	newConfig.Action(func(ctx *kingpin.ParseContext) error {
		fmt.Printf("Creating new configuration file: \"%s\"\n", *filePtr)

		conf := &config.Session{
			Name:          *namePtr,
			WindowScripts: []string{},
		}

		if *wCount > 0 {
			for i := 0; i < *wCount; i++ {
				w := config.Window{
					Scripts:     []string{},
					PaneScripts: []string{},
				}
				if *pCount > 0 {
					for j := 0; j < *pCount; j++ {
						w.Panes = append(w.Panes, config.Pane{Scripts: []string{}})
					}
				}

				conf.Windows = append(conf.Windows, w)
			}
		}

		f, err := os.Create(*filePtr)
		if err != nil {
			return err
		}

		defer f.Close()

		return toml.NewEncoder(f).Encode(conf)
	})
}
