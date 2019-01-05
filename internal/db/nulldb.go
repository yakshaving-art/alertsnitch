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
	log.Printf("save alert %#v\n", data)
	return nil
}

// Ping implements Storer interface
func (NullDB) Ping() error {
	log.Println("pong")
	return nil
}

// CheckModel implements Storer interface
func (NullDB) CheckModel() error {
	log.Println("check model")
	return nil
}

func (NullDB) String() string {
	return "null database driver"
}
