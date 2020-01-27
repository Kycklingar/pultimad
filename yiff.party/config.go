package yp

import (
	"github.com/kycklingar/pultimad/config"
)

type Config struct {
	Connstr string
}

func (c Config) Default() config.Config {
	return Config{}
}
