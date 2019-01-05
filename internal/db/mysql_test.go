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
	a.NotNil(driver)
	a.NoError(driver.Ping())
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
