package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kycklingar/pultimad/config"
	"github.com/kycklingar/pultimad/daemon"
	yp "github.com/kycklingar/pultimad/yiff.party"
)

type cfg struct {
	YiffParty yp.Config
	SleepTime time.Duration
}

func (c cfg) Default() config.Config {
	c.YiffParty = c.YiffParty.Default().(yp.Config)
	c.SleepTime = time.Second * 2
	return c
}

func setupYP(conf config.Config) *yp.Checker {
	var ypc = new(yp.Checker)

	err := ypc.Init(conf)
	if err != nil {
		log.Fatal(err)
	}

	return ypc
}

func main() {
	var conf cfg

	var configFile = flag.String("cfg", "config.json", "Config file")
	var configDefaults = flag.Bool("defaults", false, "Write a config file with default values")

	flag.Parse()

	if *configDefaults {
		err := config.Write(*configFile, conf.Default())
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	conf = conf.Default().(cfg)
	err := config.Load(*configFile, &conf)
	if err != nil {
		log.Fatal(err)
	}

	err = config.Write(*configFile, conf)
	if err != nil {
		log.Fatal(err)
	}

	daemon := daemon.NewDaemon()

	ypc := setupYP(conf.YiffParty)

	daemon.RegisterDomain("yiff.party", ypc, conf.SleepTime)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT)
	daemon.Loop(quit)

}
