package tmux

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Result is a commadn result
type Result struct {
	Stdout string
	Stderr string
}

// Exec runs a tmux command
func Exec(args ...string) (Result, error) {
	var stdin bytes.Buffer
	var stderr bytes.Buffer
	var stdout bytes.Buffer

	fmt.Println(args)

	cmd := exec.Command("tmux", args...)
	cmd.Stdin = &stdin
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return Result{}, fmt.Errorf("Error running command \"tmux %v\", %s", args, stderr.String())
	}

	return Result{stdout.String(), stderr.String()}, nil
}

// SessionInfo infos about a running tmux session
type SessionInfo struct{}

// ListSessions returns the list of sessions currently running
func ListSessions() (map[string]SessionInfo, error) {
	sessionMap := make(map[string]SessionInfo)

	res, err := Exec("ls")
	if err != nil {
		return sessionMap, nil
	}

	splits := strings.Split(res.Stdout, "\n")
	for _, sess := range splits {
		sessSplits := strings.Split(sess, ":")
		if len(sessSplits) > 1 {
			sessionMap[sessSplits[0]] = SessionInfo{}
		}
	}

	return sessionMap, nil
}
