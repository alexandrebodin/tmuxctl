package tmux

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
	var fnName string
	var fnArgs []string
	var in = Result{}
	ExecFunc = func(name string, args ...string) (Result, error) {
		fnName = name
		fnArgs = args
		return in, nil
	}

	res, err := Exec("ls")
	assert.Nil(t, err)
	assert.Equal(t, in, res)
	assert.Equal(t, "tmux", fnName)
	assert.Equal(t, []string{"ls"}, fnArgs)
}

func TestSendKeys(t *testing.T) {
	var fnName string
	var fnArgs []string
	ExecFunc = func(name string, args ...string) (Result, error) {
		fnName = name
		fnArgs = args
		return Result{}, nil
	}

	err := SendKeys("target", "ls")
	assert.Nil(t, err)
	assert.Equal(t, "tmux", fnName)
	assert.Equal(t, []string{"send-keys", "-R", "-t", "target", "ls", "C-m"}, fnArgs)
}

func TestSendRawKeys(t *testing.T) {
	var fnName string
	var fnArgs []string
	ExecFunc = func(name string, args ...string) (Result, error) {
		fnName = name
		fnArgs = args
		return Result{}, nil
	}

	err := SendRawKeys("target", "ls")
	assert.Nil(t, err)
	assert.Equal(t, "tmux", fnName)
	assert.Equal(t, []string{"send-keys", "-R", "-t", "target", "ls"}, fnArgs)
}
