package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/cenkalti/backoff/v4"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var Db *sql.DB

type EmployeeEntity struct {
	gorm.Model
	FirstName    string
	LastName     string
	Username     string
	Password     string
	Email        string
	DOB          string
	DepartmentID int64
	Position     string
}

type DepartmentEntity struct {
	gorm.Model
	Name string
}

func InitDB() (*sql.DB, error) {
	log.Println("Connecting to database...")

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 1 * time.Minute

	var db *gorm.DB
	var err error
	err = backoff.Retry(func() error {
		db, err = gorm.Open(sqlserver.Open("sqlserver://sa:yourStrong(!)Password@127.0.0.1:1433?database=master"), &gorm.Config{})
		if err != nil {
			return err
		}

		log.Println("Connected to database")

		// Migrate the schemas
		log.Println("Migrating schemas...")
		db.AutoMigrate(
			EmployeeEntity{},
			DepartmentEntity{},
		)

		dbCtx, err := db.DB()
		if err != nil {
			log.Panic(err)
		}

		// Ping the database to check if it's alive.
		if err := dbCtx.PingContext(context.Background()); err != nil {
			return err
		}

		Db = dbCtx
		return nil
	}, b)

	if err != nil {
		return nil, err
	}

	return Db, nil
}
