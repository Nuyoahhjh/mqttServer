package storage

import (
	"codeskserver/log"
	"strings"

	sqlcipher "github.com/jackfr0st13/gorm-sqlite-cipher"
	"gorm.io/gorm"
)

var gStorage *gorm.DB

func CreateFile(filename string) {

	if strings.HasPrefix(filename, "sqlite3cipher://") {
		dbfilename := filename[len("sqlite3cipher://"):]
		db, err := gorm.Open(sqlcipher.Open(dbfilename), &gorm.Config{})
		if err != nil {
			log.Error("Failed to create db:%s", dbfilename)
			panic(err)
		}
		gStorage = db
	} else {
		log.Error("Failed to create db")
		panic("Invalid db filename")
	}
}

func Instance() *gorm.DB {
	return gStorage
}
