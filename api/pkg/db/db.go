package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB struct points to DB instance.
type DB struct {
	Grm *gorm.DB
}

// GetDB accepts connection string, establishes connection to database & returns pointer to DB instance.
func GetDB(connStr string) (*DB, error) {
	db, err := get(connStr)

	if err != nil {
		return nil, err
	}

	return &DB{
		Grm: db,
	}, nil
}

// CloseDB returns method to close db connection.
func (d *DB) CloseDB() error {
	db, err := d.Grm.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func get(connStr string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
