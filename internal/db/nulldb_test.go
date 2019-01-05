package db_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/yakshaving.art/alertsnitch/internal/db"
	"testing"
)

func TestNullDBObject(t *testing.T) {
	a := assert.New(t)

	n := db.NullDB{}
	a.Equal(n.String(), "null database driver")

	a.Nil(n.Save(nil))
	a.NoError(n.Ping())
	a.NoError(n.CheckModel())
}
