package term

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GetDimensions returns a term dimensions
func GetDimensions() (string, string, error) {
	cmd := exec.Command("stty", "size")
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", "", fmt.Errorf("Error geting term size %v, %s", err, stderr.String())
	}
	s := strings.Split(strings.Trim(stdout.String(), " \n"), " ")
	height := s[0]
	width := s[1]

	return width, height, err
}
