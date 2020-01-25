package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kycklingar/pultimad/daemon"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Connection string required")
		os.Exit(1)
	}

	daemon, err := daemon.NewDaemon(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	//log.Fatal(daemon.ReloadCreators())
	daemon.Loop()
}
