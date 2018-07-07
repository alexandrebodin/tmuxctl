package main

import (
	"github.com/alexandrebodin/tmuxctl/cli"
)

var (
	version = "development"
)

func main() {
	cli.Run(version)
}
