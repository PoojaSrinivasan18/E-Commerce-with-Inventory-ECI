package database

import (
	"fmt"
	"time"

	"github.com/PoojaSrinivasan18/payment-service/common"
	"github.com/PoojaSrinivasan18/payment-service/model"

	"github.com/apex/log"
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
	log.Infof("entering setupDb")
	var db *gorm.DB

	driver := configuration.Database.Driver
	//dbname := configuration.Database.Dbname
	//username := configuration.Database.Username
	password := configuration.Database.Password

	//port := configuration.Database.Port

	//host := os.Getenv("MY_POD_IP")
	host := configuration.Database.Host
	if host != "" {
		log.Infof("Host IP is %v", host)
	} else {
		log.Error("Host is Empty in Env Variable")
	}

	// data source name
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		configuration.Database.Username,
		password,
		configuration.Database.Dbname,
		configuration.Database.Port,
	)

	if driver == "postgres" { // Postgres DB
		for i := 0; i < 10; i++ {
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err == nil {
				log.Infof("Successfully connected to DB on attempt %d", i+1)
				break
			}
			log.Warnf("Failed to connect to DB (attempt %d/10): %v", i+1, err)
			time.Sleep(5 * time.Second)
		}

		if err != nil {
			log.Fatalf("Could not connect to database after 10 attempts: %v", err)
		}
	}

	// Change this to true if you want to see SQL queries
	database, dbErr := db.DB()
	if dbErr != nil {
		log.Errorf("db err: ", err)
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
	err = Repo.Database.AutoMigrate(&model.PaymentModel{})
	if err != nil {
		log.Errorf("Auto-migrate error: ", err)
	}
}

func GetDB() *gorm.DB {
	return Repo.Database
}
