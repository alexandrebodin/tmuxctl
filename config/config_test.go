package config_test

import (
	"bytes"
	"testing"

	"github.com/alexandrebodin/tmuxctl/config"
	"github.com/stretchr/testify/assert"
)

type ReaderString struct {
	*bytes.Reader
}

func (r ReaderString) Close() error {
	return nil
}

func TestParse(t *testing.T) {

	t.Run("Error invalid toml", func(t *testing.T) {
		reader := ReaderString{bytes.NewReader([]byte("lala"))}
		_, err := config.Parse(reader)
		assert.Error(t, err)
	})

	t.Run("Invalid selected-window", func(t *testing.T) {
		reader := ReaderString{bytes.NewReader([]byte("select-window=\"unknwon-window\""))}
		_, err := config.Parse(reader)
		assert.Error(t, err)
	})

	t.Run("Invalid selected-pane", func(t *testing.T) {
		reader := ReaderString{bytes.NewReader([]byte(`
			select-window="win"
			select-pane=99
			[[windows]]
				name="win"
		`))}

		_, err := config.Parse(reader)
		assert.Error(t, err)
	})

	t.Run("Invalid layout", func(t *testing.T) {
		reader := ReaderString{bytes.NewReader([]byte(`
			[[windows]]
				layout="unknown"
		`))}

		_, err := config.Parse(reader)
		assert.Error(t, err)
	})

	t.Run("Valid configuration", func(t *testing.T) {
		reader := ReaderString{bytes.NewReader([]byte(`
			name="valid session"
			[[windows]]
			layout="tiled"
		`))}

		_, err := config.Parse(reader)
		assert.Nil(t, err)
	})
}
