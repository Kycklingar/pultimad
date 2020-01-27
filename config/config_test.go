package config

import (
	"strings"
	"testing"
)

func expected(t *testing.T, exp, got string) {
	t.Fatalf("Expected %s, got %s", exp, got)
}

type cconfig struct {
	Hello   string
	Darling string
}

func (c cconfig) Default() Config {
	return c
}

func TestCustomConfig(t *testing.T) {
	var cc cconfig
	cc.Hello = "Worldstar"
	cc.Darling = "Merry"

	ss := strings.NewReader("{\"Hello\":\"World\"}")
	err := Parse(ss, &cc)
	if err != nil {
		t.Fatal(err)
	}

	if cc.Hello != "World" {
		expected(t, "World", cc.Hello)
	}

	if cc.Darling != "Merry" {
		expected(t, "Merry", cc.Darling)
	}
}

func TestLoad(t *testing.T) {
	var cfg cconfig
	err := Load("test.cfg", &cfg)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Hello != "World!" {
		expected(t, "World", cfg.Hello)
	}
}

func TestWrite(t *testing.T) {
	var cfg cconfig
	cfg.Hello = "Oxy"
	cfg.Darling = "Daisy"

	if err := Write("testWrite.cfg", cfg); err != nil {
		t.Fatal(err)
	}
}
