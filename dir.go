package main

import (
	"log"
	"os/user"
	"path/filepath"
	"strings"
)

func lookupDir(path string) string {
	if path == "" {
		return path
	}

	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}

		if usr.HomeDir != "" {
			return filepath.Join(usr.HomeDir, path[1:])
		}

		log.Fatal("No home directory")
	}

	return path
}
