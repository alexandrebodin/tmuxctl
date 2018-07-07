package cli

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
)

var rootCmd = kingpin.New("tmuxctl", "")

// Run run the tmuxctl cli
func Run(version string) {
	rootCmd.Version(fmt.Sprintf("tmuxct %s", version)).Author("Alexandre BODIN")
	rootCmd.HelpFlag.Short('h')
	rootCmd.VersionFlag.Short('v')
	kingpin.MustParse(rootCmd.Parse(os.Args[1:]))
}
