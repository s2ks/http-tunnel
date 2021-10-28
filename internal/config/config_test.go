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

	if conf["__global"]["this"][0] != "exists in [__global]" {
		t.Fail()
	}
	if conf["test_section1"]["key1"][0] != "foo1" {
		t.Fail()
	}
	if conf["test_section1"]["key1"][1] != "foo1.1" {
		t.Fail()
	}
	if conf["test_section1"]["key2"][0] != "bar1" {
		t.Fail()
	}
	if conf["test_section1"]["this"][0] != "not a comment" {
		t.Fail()
	}
	if len(conf["empty"]) != 0 {
		t.Fail()
	}
        if conf["not empty"]["keys can have spaces"][0] != "" {
                t.Fail()
        }
        if conf["not empty"]["and don't need"][0] != "a val (default empty)" {
                t.Fail()
        }

	t.Log(conf)
}
