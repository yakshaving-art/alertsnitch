package db

import (
	"log"

	"github.com/prometheus/alertmanager/template"
)

// NullDB A database that does nothing
type NullDB struct {
}

// Save implements Storer interface
func (NullDB) Save(data *template.Data) error {
	log.Printf("save alert %#v", data)
	return nil
}

// Ping implements Storer interface
func (NullDB) Ping() error {
	log.Printf("pong")
	return nil
}

func (NullDB) String() string {
	return "null database driver"
}
