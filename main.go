package main

import (
	"flag"
	"fmt"
	"os"

	"gitlab.com/yakshaving.art/alertsnitch/internal/db"
	"gitlab.com/yakshaving.art/alertsnitch/internal/server"
	"gitlab.com/yakshaving.art/alertsnitch/version"
)

// Args are the arguments that can be passed to alertsnitch
type Args struct {
	Version bool
	Address string
}

func main() {
	args := Args{}

	flag.BoolVar(&args.Version, "version", false, "print the version and exit")
	flag.StringVar(&args.Address, "listen.address", ":8080", "address in which to listen for http requests")

	flag.Parse()

	if args.Version {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	}

	s := server.New(db.NullDB{})
	s.Start(args.Address)
}
