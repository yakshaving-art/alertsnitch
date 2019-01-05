package internal

import (
	"github.com/prometheus/alertmanager/template"
)

// Storer saves an Alert Data into a persistence engine
type Storer interface {
	Save(*template.Data) error
	Ping() error
	CheckModel() error
}
