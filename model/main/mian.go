package main

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	viper.SetConfigName("settings")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	username := viper.Get("settings.mysql.username").(string)
	password := viper.Get("settings.mysql.password").(string)
	host := viper.Get("settings.mysql.host").(string)
	port := viper.Get("settings.mysql.port").(int)
	schema := viper.Get("settings.mysql.schema").(string)
	parameters := viper.Get("settings.mysql.parameters").(string)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", username, password, host, port, schema, parameters)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	err = db.AutoMigrate(
		&model.User{},
		&model.Video{},
		&model.Comment{},
		&model.Like{},
		&model.Message{},
		&model.Relation{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
