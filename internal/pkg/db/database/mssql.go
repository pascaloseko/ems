package database

import (
	"context"
	"database/sql"
	"log"

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
	db, err := gorm.Open(sqlserver.Open("sqlserver://sa:yourStrong(!)Password@localhost:1433?database=master"), &gorm.Config{})
	if err != nil {
		return nil, err
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
		return nil, err
	}

	Db = dbCtx
	return dbCtx, nil
}
