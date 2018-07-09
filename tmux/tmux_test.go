package tmux

import (
	"errors"
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

func TestGetOptions(t *testing.T) {
	t.Run("get options", func(t *testing.T) {
		ExecFunc = func(name string, args ...string) (Result, error) {
			return Result{Stdout: "base-index 2\npane-base-index 1"}, nil
		}

		options, err := GetOptions()
		assert.Nil(t, err)
		assert.Equal(t, 2, options.BaseIndex)
		assert.Equal(t, 1, options.PaneBaseIndex)
	})

	t.Run("get options error", func(t *testing.T) {
		ExecFunc = func(name string, args ...string) (Result, error) {
			return Result{}, errors.New("some random error")
		}

		_, err := GetOptions()
		assert.Error(t, err)
	})
}

func TestListSessions(t *testing.T) {
	t.Run("list options ok", func(t *testing.T) {
		ExecFunc = func(name string, args ...string) (Result, error) {
			return Result{Stdout: "azd: infos\nother-session: infos"}, nil
		}

		sessions, err := ListSessions()
		assert.Nil(t, err)
		assert.Equal(t, map[string]SessionInfo{
			"azd":           SessionInfo{},
			"other-session": SessionInfo{},
		},
			sessions,
		)
	})

	t.Run("list options error", func(t *testing.T) {
		ExecFunc = func(name string, args ...string) (Result, error) {
			return Result{}, errors.New("some error")
		}

		sessions, err := ListSessions()
		assert.Nil(t, err)
		assert.Equal(t, map[string]SessionInfo{}, sessions)
	})
}
