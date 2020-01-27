package main

import (
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

func main() {
	var conf = new(cfg)

	if len(os.Args) > 1 {
		err := config.Write("config.json", conf.Default())
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	err := config.Load("config.json", conf)
	if err != nil {
		log.Fatal(err)
		return
	}

	var ypc = new(yp.Checker)

	err = ypc.Init(conf.YiffParty)
	if err != nil {
		log.Fatal(err)
		return
	}

	daemon := daemon.NewDaemon()

	daemon.RegisterDomain("yiff.party", ypc, conf.SleepTime)

	//log.Fatal(daemon.ReloadCreators())

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGINT)

	daemon.Loop(quit)

}
