package db

import (
	"context"
	"fmt"
	"time"

	"database/sql"

	"gitlab.com/yakshaving.art/alertsnitch/internal"
	"gitlab.com/yakshaving.art/alertsnitch/internal/metrics"
)

// SupportedModel stores the model that is supported by this application
const SupportedModel = "0.0.1"

// MySQLDB A database that does nothing
type MySQLDB struct {
	db *sql.DB
}

// ConnectMySQL connect to a MySQL database using the provided data source name
func ConnectMySQL(dsn string) (*MySQLDB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("Empty DSN provided, can't connect to database")
	}

	connection, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL connection: %s", err)
	}

	database := &MySQLDB{
		db: connection,
	}

	err = database.Ping()
	if err != nil {
		return nil, err
	}

	return database, database.CheckModel()
}

// Save implements Storer interface
func (d MySQLDB) Save(data *internal.AlertGroup) error {
	return d.unitOfWork(func() error {
		r, err := d.db.Exec(`
			INSERT INTO AlertGroup (time, receiver, status, externalURL, groupKey)
			VALUES (now(), ?, ?, ?, ?)`, data.Receiver, data.Status, data.ExternalURL, data.GroupKey)
		if err != nil {
			return fmt.Errorf("failed to insert into AlertGroups: %s", err)
		}

		alertGroupID, err := r.LastInsertId() // alertGroupID
		if err != nil {
			return fmt.Errorf("failed to get AlertGroups inserted id: %s", err)
		}

		for k, v := range data.GroupLabels {
			_, err := d.db.Exec(`
				INSERT INTO GroupLabel (alertGroupID, GroupLabel, Value)
				VALUES (?, ?, ?)`, alertGroupID, k, v)
			if err != nil {
				return fmt.Errorf("failed to insert into GroupLabel: %s", err)
			}
		}
		for k, v := range data.CommonLabels {
			_, err := d.db.Exec(`
				INSERT INTO CommonLabel (alertGroupID, Label, Value)
				VALUES (?, ?, ?)`, alertGroupID, k, v)
			if err != nil {
				return fmt.Errorf("failed to insert into CommonLabel: %s", err)
			}
		}
		for k, v := range data.CommonAnnotations {
			_, err := d.db.Exec(`
				INSERT INTO CommonAnnotation (alertGroupID, Annotation, Value)
				VALUES (?, ?, ?)`, alertGroupID, k, v)
			if err != nil {
				return fmt.Errorf("failed to insert into CommonAnnotation: %s", err)
			}
		}

		for _, alert := range data.Alerts {
			var result sql.Result
			if alert.EndsAt.Before(alert.StartsAt) {
				result, err = d.db.Exec(`
				INSERT INTO Alert (alertGroupID, status, startsAt, generatorURL)
				VALUES (?, ?, ?, ?)`,
					alertGroupID, alert.Status, alert.StartsAt, alert.EndsAt, alert.GeneratorURL)
			} else {
				result, err = d.db.Exec(`
				INSERT INTO Alert (alertGroupID, status, startsAt, endsAt, generatorURL)
				VALUES (?, ?, ?, ?, ?)`,
					alertGroupID, alert.Status, alert.StartsAt, alert.EndsAt, alert.GeneratorURL)
			}
			if err != nil {
				return fmt.Errorf("failed to insert into Alert: %s", err)
			}

			alertID, err := result.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get Alert inserted id: %s", err)
			}

			for k, v := range alert.Labels {
				_, err := d.db.Exec(`
					INSERT INTO AlertLabel (AlertID, Label, Value)
					VALUES (?, ?, ?)`, alertID, k, v)
				if err != nil {
					return fmt.Errorf("failed to insert into AlertLabel: %s", err)
				}
			}
			for k, v := range alert.Annotations {
				_, err := d.db.Exec(`
					INSERT INTO AlertAnnotation (AlertID, Annotation, Value)
					VALUES (?, ?, ?)`, alertID, k, v)
				if err != nil {
					return fmt.Errorf("failed to insert into AlertAnnotation: %s", err)
				}
			}
		}

		return nil
	})
}

func (d MySQLDB) unitOfWork(f func() error) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %s", err)
	}

	err = f()

	if err != nil {
		e := tx.Rollback()
		if e != nil {
			return fmt.Errorf("failed to rollback transaction (%s) after failing execution: %s", e, err)
		}
		return fmt.Errorf("failed execution: %s", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %s", err)
	}
	return nil
}

// Ping implements Storer interface
func (d MySQLDB) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := d.db.PingContext(ctx); err != nil {
		metrics.DatabaseUp.Set(0)
		return err
	}
	metrics.DatabaseUp.Set(1)
	return nil
}

// CheckModel implements Storer interface
func (d MySQLDB) CheckModel() error {
	rows, err := d.db.Query("SELECT version FROM Model")
	if err != nil {
		return fmt.Errorf("failed to fetch model version from the database: %s", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return fmt.Errorf("failed to read model version from the database: empty resultset")
	}

	var model string
	if err := rows.Scan(&model); err != nil {
		return fmt.Errorf("failed to read model version from the database: %s", err)
	}

	if model != SupportedModel {
		return fmt.Errorf("model '%s' is not supported by this application (%s)", model, SupportedModel)
	}

	return nil
}

func (MySQLDB) String() string {
	return "mysql database driver"
}
