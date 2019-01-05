// +build integration

package db_test

import (
	"io/ioutil"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"

	"gitlab.com/yakshaving.art/alertsnitch/internal/db"
	"gitlab.com/yakshaving.art/alertsnitch/internal/webhook"
)

func TestPingingDatabaseWorks(t *testing.T) {
	a := assert.New(t)
	driver, err := db.ConnectMySQL(os.Getenv("ALERTSNITCHER_MYSQL_DSN"))
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

	driver, err := db.ConnectMySQL(os.Getenv("ALERTSNITCHER_MYSQL_DSN"))
	a.NoError(err)

	a.NoError(driver.Save(data))
}
