package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kycklingar/pultimad/daemon"
	yp "github.com/kycklingar/pultimad/yiff.party"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Connection string required")
		os.Exit(1)
	}

	var ypc = new(yp.Checker)

	ypc.Init(os.Args[1])

	daemon := daemon.NewDaemon()

	daemon.RegisterDomain("yiff.party", ypc, time.Second*2)

	//log.Fatal(daemon.ReloadCreators())

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGINT)

	daemon.Loop(quit)

}
