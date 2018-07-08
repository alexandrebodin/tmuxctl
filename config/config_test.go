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
	reader := ReaderString{bytes.NewReader([]byte("lala"))}

	_, err := config.Parse(reader)
	assert.Error(t, err)
}
