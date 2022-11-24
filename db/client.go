package db

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *DB

type DB struct {
	*gorm.DB
}

func GetDB() (*DB, error) {
	if db != nil {
		return db, nil
	}
	_db, err := gorm.Open(postgres.Open("host=localhost user=sarah password=password dbname=skynet port=5432 sslmode=disable TimeZone=UTC"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	db = &DB{_db}
	return db, nil
}

func (db *DB) GetSQLConn() (*sql.DB, error) {
	return db.DB.DB()
}
