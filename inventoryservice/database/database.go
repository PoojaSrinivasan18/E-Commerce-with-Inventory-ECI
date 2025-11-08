package database

import (
	common "inventoryservice/common"
	models "inventoryservice/models"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Repo = Repository{}
	err  error
)

type Repository struct {
	Database *gorm.DB
}

// SetupDB opens a database and saves the reference to `Database` struct.
func SetupDB(configuration *common.Configuration) error {
	var db *gorm.DB

	driver := configuration.Database.Driver
	dbname := configuration.Database.Dbname
	username := configuration.Database.Username
	password := configuration.Database.Password

	port := configuration.Database.Port

	//host := os.Getenv("MY_POD_IP")
	host := configuration.Database.Host
	if host != "" {
		log.Info("Host IP is ", host)
	} else {
		log.Error("Host is Empty in Env Variable")
	}

	pw := os.Getenv("APP_DB_PASSWORD")
	if pw != "" {
		password = pw
	}

	// data source name
	dsn := "host=" + host + " user=" + username + " password=" + password + " port=" + port + " dbname=" + dbname
	if driver == "postgres" { // Postgres DB
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Error("db err: ", err)
		}
	}

	// Change this to true if you want to see SQL queries
	database, dbErr := db.DB()
	if dbErr != nil {
		log.Error("db err: ", err)
		return err
	}
	database.SetMaxIdleConns(configuration.Database.MaxIdleConns)
	database.SetMaxOpenConns(configuration.Database.MaxOpenConns)
	database.SetConnMaxLifetime(time.Duration(configuration.Database.MaxLifetime) * time.Second)
	Repo.Database = db
	migrateModels()

	return nil
}

// Auto migrate project models
func migrateModels() {
	err = Repo.Database.AutoMigrate(&models.InventoryModel{}, &models.ReservationRecord{})
	if err != nil {
		log.Error("Auto-migrate error: ", err)
	}
}

func GetDB() *gorm.DB {
	return Repo.Database
}
