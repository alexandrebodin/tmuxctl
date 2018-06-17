package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/BurntSushi/toml"
)

type Session struct {
	Name string
	Dir  string
}

func (sess *Session) Sart() error {
	cmd := exec.Command("tmux", "new-session", "-d", "-s", sess.Name)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error Creating tmux session %s", err)
	}

	tmux, err := exec.LookPath("tmux")
	if err != nil {
		return fmt.Errorf("Error looking up tmux %s", err)
	}

	args := []string{"tmux", "attach", "-t", sess.Name}
	if sysErr := syscall.Exec(tmux, args, os.Environ()); sysErr != nil {
		return fmt.Errorf("Error attaching to session %s, %s", sess.Name, sysErr)
	}
	return nil
}

func main() {
	args := []string{".tmuxrc"}

	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	filePath := args[0]

	var conf map[string]interface{}
	if _, err := toml.DecodeFile(filePath, &conf); err != nil {
		panic(fmt.Errorf("Error decoding configuration %s", err))
	}

	sess := &Session{
		Name: conf["name"].(string),
		Dir:  conf["cwd"].(string),
	}

	err := sess.Sart()

	if err != nil {
		fmt.Println(err)
	}
}
