/*
 * @File: config.database.go
 * @Description: Defines database information of the service
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

  var DB *gorm.DB

  func ConfigDatabase(){
    dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
        os.Getenv("POSTGRES_HOST"),
        os.Getenv("POSTGRES_USER"),
        os.Getenv("POSTGRES_PASSWORD"),
        os.Getenv("POSTGRES_DB"),
        os.Getenv("POSTGRES_PORT"),
        os.Getenv("POSTGRES_TIMEZONE"),
    )
    db, err := gorm.Open(postgres.New(postgres.Config{
        DSN:                  dns,
        PreferSimpleProtocol: true,
    }), &gorm.Config{})
    if err != nil {
        return
    }
    fmt.Print(db)
    DB = db
}
