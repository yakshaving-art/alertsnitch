package db

import (
	"github.com/prometheus/alertmanager/template"
)

// NullDB A database that does nothing
type NullDB struct {
}

// Save implements Storer interface
func (NullDB) Save(data *template.Data) error {
	return nil
}

// Ping implements Storer interface
func (NullDB) Ping() error {
	return nil
}

func (NullDB) String() string {
	return "null database driver"
}
