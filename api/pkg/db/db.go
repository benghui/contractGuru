package db

import (
	"encoding/base64"
	"os"

	"github.com/wader/gormstore/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB struct points to gorm DB & Store instance.
type DB struct {
	Grm   *gorm.DB
	Store *gormstore.Store
}

// GetDB accepts connection string, establishes connection to database & returns pointer to DB instance.
func GetDB(connStr string) (*DB, error) {
	db, err := get(connStr)

	if err != nil {
		return nil, err
	}

	authKey := encode([]byte(os.Getenv("AKEY")))
	encryptionKey := encode([]byte(os.Getenv("EKEY")))

	store := gormstore.NewOptions(db, gormstore.Options{}, authKey, encryptionKey)

	store.SessionOpts.HttpOnly = true
	store.SessionOpts.Secure = true
	store.SessionOpts.MaxAge = 86400 * 7

	return &DB{
		Grm:   db,
		Store: store,
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

func encode(value []byte) []byte {
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(value)))
	base64.URLEncoding.Encode(encoded, value)
	return encoded
}
