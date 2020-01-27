package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

var (
	ErrInvalid = errInvalid()
)

func errInvalid() error {
	return errors.New("Invalid config type")
}

type Config interface {
	Default() Config
}

func Load(filename string, out Config) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return Parse(f, out)
}

func Parse(input io.Reader, c Config) error {
	dec := json.NewDecoder(input)
	err := dec.Decode(&c)
	if err != nil {
		return err
	}

	return nil
}

func Write(filename string, cfg Config) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "	")
	return enc.Encode(cfg)
}
