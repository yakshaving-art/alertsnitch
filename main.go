package main

import (
	"flag"
	"fmt"
	"os"

	"gitlab.com/yakshaving.art/alertsnitch/version"
)

// Args are the arguments that can be passed to alertsnitch
type Args struct {
	Version bool
}

func main() {
	args := Args{}

	flag.BoolVar(&args.Version, "version", false, "print the version and exit")

	flag.Parse()

	if args.Version {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	}

}
