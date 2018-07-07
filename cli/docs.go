package cli

import (
	"github.com/alecthomas/kingpin"
	"github.com/pkg/browser"
)

var url = "https://tmuxctl.netlify.com"

func init() {
	var docs = rootCmd.Command("docs", "Open documentation website")

	docs.Action(func(ctx *kingpin.ParseContext) error {
		browser.OpenURL(url)
		return nil
	})
}
