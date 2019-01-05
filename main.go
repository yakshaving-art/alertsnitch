package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"gitlab.com/yakshaving.art/alertsnitch/internal/db"
	"gitlab.com/yakshaving.art/alertsnitch/internal/server"
	"gitlab.com/yakshaving.art/alertsnitch/version"
)

// Args are the arguments that can be passed to alertsnitch
type Args struct {
	Version bool
	Address string
	DSN     string
}

func main() {
	args := Args{
		DSN: os.Getenv("ALERTSNITCHER_MYSQL_DSN"),
	}

	flag.BoolVar(&args.Version, "version", false, "print the version and exit")
	flag.StringVar(&args.Address, "listen.address", ":8080", "address in which to listen for http requests")

	flag.Parse()

	if args.Version {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	}

	driver, err := db.ConnectMySQL(args.DSN)
	if err != nil {
		fmt.Println("failed to connect to MySQL database: ", err)
		os.Exit(1)
	}

	s := server.New(driver)
	s.Start(args.Address)
}
