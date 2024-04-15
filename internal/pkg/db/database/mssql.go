package database

import (
	"database/sql"
	"fmt"
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

var server = "sqlserver"
var port = 1433
var databaseName = "mydb"

func InitDB() {
	connString := fmt.Sprintf("server=%s;port=%d;database=%s;fedauth=ActiveDirectoryDefault;", server, port, databaseName)
	db, err := gorm.Open(sqlserver.Open(connString), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	db.AutoMigrate(
		EmployeeEntity{},
		DepartmentEntity{},
	)

	dbCtx, err := db.DB()
	if err != nil {
		log.Panic(err)
	}

	defer dbCtx.Close()
	Db = dbCtx
}

func CloseDB() error {
	return Db.Close()
}
