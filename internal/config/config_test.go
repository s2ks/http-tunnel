package config

import (
	"testing"
	"os"
)

func TestParse(t *testing.T) {
	file, err := os.Open("testdata/config.ini")

	if err != nil {
		t.Error(err)
	}

	conf, err := Parse(file)

	if err != nil {
		t.Error(err)
	}

	t.Log(conf)
}
