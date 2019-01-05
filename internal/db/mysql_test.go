// +build integration

package db_test

import (
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"

	"gitlab.com/yakshaving.art/alertsnitch/internal/db"
)

func TestIntegration(t *testing.T) {
	a := assert.New(t)
	driver, err := db.ConnectMySQL(os.Getenv("ALERTSNITCHER_MYSQL_DSN"))
	a.NoError(err)
	a.NotNil(driver)
	a.NoError(driver.Ping())
}
