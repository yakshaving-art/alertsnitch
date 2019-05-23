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
	Address                string
	DSN                    string
	MaxIdleConns           int
	MaxOpenConns           int
	MaxConnLifetimeSeconds int

	Debug   bool
	DryRun  bool
	Version bool
}

func main() {
	args := Args{
		DSN: os.Getenv(internal.MySQLDSNVar),
	}

	flag.BoolVar(&args.Version, "version", false, "print the version and exit")
	flag.StringVar(&args.Address, "listen.address", ":9567", "address in which to listen for http requests")
	flag.BoolVar(&args.DryRun, "dryrun", false, "uses a null db driver that writes received webhooks to stdout")
	flag.BoolVar(&args.Debug, "debug", false, "enable debug mode, which dumps alerts payloads to the log as they arrive")

	flag.IntVar(&args.MaxOpenConns, "max-open-connections", 2, "maximum number of connections in the pool")
	flag.IntVar(&args.MaxIdleConns, "max-idle-connections", 1, "maximum number of idle connections in the pool")
	flag.IntVar(&args.MaxConnLifetimeSeconds, "max-connection-lifetyme-seconds", 600, "maximum number of seconds a connection is kept alive in the pool")

	flag.Parse()

	if args.Version {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	}

	var driver internal.Storer
	if args.DryRun {
		driver = db.NullDB{}

	} else {
		d, err := db.ConnectMySQL(db.ConnectionArgs{
			DSN:                    args.DSN,
			MaxIdleConns:           args.MaxIdleConns,
			MaxOpenConns:           args.MaxOpenConns,
			MaxConnLifetimeSeconds: args.MaxConnLifetimeSeconds,
		})
		if err != nil {
			fmt.Println("failed to connect to MySQL database: ", err)
			os.Exit(1)
		}
		driver = d
	}

	s := server.New(driver, args.Debug)
	s.Start(args.Address)
}
