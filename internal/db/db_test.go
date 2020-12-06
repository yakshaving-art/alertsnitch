// +build integration

package db_test

import (
	"io/ioutil"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"

	"gitlab.com/yakshaving.art/alertsnitch/internal"
	"gitlab.com/yakshaving.art/alertsnitch/internal/db"
	"gitlab.com/yakshaving.art/alertsnitch/internal/webhook"
)

func TestPingingDatabaseWorks(t *testing.T) {
	backend := os.Getenv("ALERTSNITCH_BACKEND")

	a := assert.New(t)
	driver, err := db.Connect(backend, connectionArgs())
	a.NoError(err)
	a.NotNilf(driver, "database driver is nil?")
	a.NoErrorf(driver.Ping(), "failed to ping database")
	a.NoErrorf(driver.CheckModel(), "failed to check the model")
}

func TestSavingAnAlertWorks(t *testing.T) {
	a := assert.New(t)

	b, err := ioutil.ReadFile("../webhook/sample-payload.json")
	a.NoError(err)

	data, err := webhook.Parse(b)
	a.NoError(err)

	backend := os.Getenv("ALERTSNITCH_BACKEND")

	driver, err := db.Connect(backend, connectionArgs())
	a.NoError(err)

	a.NoError(driver.Save(data))
}

func TestSavingAFiringAlertWorks(t *testing.T) {
	a := assert.New(t)

	b, err := ioutil.ReadFile("../webhook/sample-payload-invalid-ends-at.json")
	a.NoError(err)

	data, err := webhook.Parse(b)
	a.NoError(err)

	backend := os.Getenv("ALERTSNITCH_BACKEND")
	driver, err := db.Connect(backend, connectionArgs())
	a.NoError(err)

	a.NoError(driver.Save(data))
}

func connectionArgs() db.ConnectionArgs {
	return db.ConnectionArgs{
		DSN:                    os.Getenv(internal.DSNVar),
		MaxIdleConns:           1,
		MaxOpenConns:           2,
		MaxConnLifetimeSeconds: 600,
	}
}
