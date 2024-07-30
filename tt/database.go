package tt

import (
	"errors"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {

	dbfile := viper.GetString("db")
	if dbfile == "" {
		log.Fatalf(ErrorString, CharError, errors.New("please `export TT_DB` to the location the zeit database should be stored at or create a config file"))
	}

	var err error
	db, err = gorm.Open(sqlite.Open(viper.GetString("db")), &gorm.Config{})
	if err != nil {
		log.Fatalf(ErrorString, CharError, err)
	}

	clientInit()
	projectInit()
}
