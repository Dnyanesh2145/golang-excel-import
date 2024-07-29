package dialects

import (
	"errors"
	"golang-excel-import/config"
	"golang-excel-import/models"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var once sync.Once
var DBconnection *gorm.DB

func GetConnection() (*gorm.DB, error) {
	once.Do(func() {
		var err error
		dbURI := config.GetLocalEnv("DB_URI")
		DBconnection, err = connect(dbURI)
		if err != nil {
			log.Panic("Error while connecting to Mysql DB")
			log.Panic(err)
			return
		}
	})
	return DBconnection, nil
}

func connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		log.Panic(err)
		return nil, errors.New("Cannot connect to postgres")
	}
	db.Logger.LogMode(logger.LogLevel(logger.Info))
	log.Println("connected to postgress DB.")
	return db, nil
}

func Close(db *gorm.DB) {
	conn, err := db.DB()
	if err != nil {
		log.Panic("Error while closing this error")
	}
	Conn, _ := DBconnection.DB()
	Conn.Close()
	conn.Close()
}

func Ping() error {
	if val, err := DBconnection.DB(); err != nil {
		return err
	} else {
		return val.Ping()
	}
}

func CheckModals() {
	conn, err := GetConnection()
	if err != nil {
		log.Panic("Error while running Postgress Migration")
		return
	}

	conn.AutoMigrate(&models.Records{})
	

}
