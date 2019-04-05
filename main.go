package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"gitlab.com/yakshaving.art/alertsnitch/internal"
	"gitlab.com/yakshaving.art/alertsnitch/internal/db"
	"gitlab.com/yakshaving.art/alertsnitch/internal/server"
	"gitlab.com/yakshaving.art/alertsnitch/version"
)

// Args are the arguments that can be passed to alertsnitch
type Args struct {
	Address string
	DSN     string

	Debug   bool
	DryRun  bool
	Version bool
}

func main() {
	args := Args{
		DSN: os.Getenv(internal.MySQLDSNVar),
	}

	flag.BoolVar(&args.Version, "version", false, "print the version and exit")
	flag.StringVar(&args.Address, "listen.address", ":8080", "address in which to listen for http requests")
	flag.BoolVar(&args.DryRun, "dryrun", false, "uses a null db driver that writes received webhooks to stdout")
	flag.BoolVar(&args.Debug, "debug", false, "enable debug mode, which dumps alerts payloads to the log as they arrive")

	flag.Parse()

	if args.Version {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	}

	var driver internal.Storer
	if args.DryRun {
		driver = db.NullDB{}

	} else {
		d, err := db.ConnectMySQL(args.DSN)
		if err != nil {
			fmt.Println("failed to connect to MySQL database: ", err)
			os.Exit(1)
		}
		driver = d
	}

	s := server.New(driver, args.Debug)
	s.Start(args.Address)
}
