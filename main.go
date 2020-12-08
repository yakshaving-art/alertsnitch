package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/sirupsen/logrus"

	"gitlab.com/yakshaving.art/alertsnitch/internal"
	"gitlab.com/yakshaving.art/alertsnitch/internal/db"
	"gitlab.com/yakshaving.art/alertsnitch/internal/server"
	"gitlab.com/yakshaving.art/alertsnitch/version"
)

// Args are the arguments that can be passed to alertsnitch
type Args struct {
	Address                string
	DBBackend              string
	DSN                    string
	MaxIdleConns           int
	MaxOpenConns           int
	MaxConnLifetimeSeconds int

	Debug   bool
	DryRun  bool
	Version bool
}

func main() {
	args := Args{}

	flag.BoolVar(&args.Version, "version", false, "print the version and exit")
	flag.StringVar(&args.Address, "listen.address", envWithDefault("ALERTSNITCH_ADDR", ":9567"), "address in which to listen for http requests, (also ALERTSNITCH_ADDR)")
	flag.BoolVar(&args.Debug, "debug", false, "enable debug mode, which dumps alerts payloads to the log as they arrive")

	flag.StringVar(&args.DBBackend, "database-backend", envWithDefault("ALERTSNITCH_BACKEND", "mysql"), "database backend, allowed are mysql, postgres, and null (also ALERTSNITCH_BACKEND")
	flag.StringVar(&args.DSN, "dsn", os.Getenv(internal.DSNVar), "Database DSN (also ALERTSNITCH_DSN)")

	flag.IntVar(&args.MaxOpenConns, "max-open-connections", 2, "maximum number of connections in the pool")
	flag.IntVar(&args.MaxIdleConns, "max-idle-connections", 1, "maximum number of idle connections in the pool")
	flag.IntVar(&args.MaxConnLifetimeSeconds, "max-connection-lifetyme-seconds", 600, "maximum number of seconds a connection is kept alive in the pool")

	flag.Parse()

	if args.Version {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	}

	if args.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	driver, err := db.Connect(args.DBBackend, db.ConnectionArgs{
		DSN:                    args.DSN,
		MaxIdleConns:           args.MaxIdleConns,
		MaxOpenConns:           args.MaxOpenConns,
		MaxConnLifetimeSeconds: args.MaxConnLifetimeSeconds,
	})
	if err != nil {
		fmt.Println("failed to connect to database:", err)
		os.Exit(1)
	}

	s := server.New(driver, args.Debug)
	s.Start(args.Address)
}

func envWithDefault(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}
