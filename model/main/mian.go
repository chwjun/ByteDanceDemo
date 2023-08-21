package main

import (
	"bytedancedemo/model"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*创建表并插入假数据*/
func main() {
	// viper.SetConfigName("settings")
	// viper.SetConfigType("yml")
	// viper.AddConfigPath(".")
	// if err := viper.ReadInConfig(); err != nil {
	// 	log.Fatalf("Error reading config file, %s", err)
	// }

	// username := viper.Get("settings.mysql.username").(string)
	// password := viper.Get("settings.mysql.password").(string)
	// host := viper.Get("settings.mysql.host").(string)
	// port := viper.Get("settings.mysql.port").(int)
	// schema := viper.Get("settings.mysql.schema").(string)
	// parameters := viper.Get("settings.mysql.parameters").(string)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", "sample_douyin", "sample_douyin", "43.140.203.85", 3306, "sample_douyin", "charset=utf8mb4&parseTime=True&loc=Local")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// models := []interface{}{
	// 	&model.User{},
	// 	&model.Video{},
	// 	&model.Comment{},
	// 	&model.Like{},
	// 	&model.Message{},
	// 	&model.Relation{},
	// }

	// for _, model := range models {
	// 	if db.Migrator().HasTable(model) {
	// 		insertFakeData(db)
	// 		log.Fatalf("Table %T already exists", model)
	// 	} else {
	// 		err = db.AutoMigrate(model)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		} else {
	// 			fmt.Printf("AutoMigrate success for %T!\n", model)
	// 		}
	// 	}
	// }
	insertFakeData(db)
}

func insertFakeData(db *gorm.DB) {
	// Insert fake data into each table

	// Insert User fake data
	//db.Create(&model.User{Name: "Alice", Password: "password123"})
	//db.Create(&model.User{Name: "Bob", Password: "password456"})

	// Insert Video fake data
	db.Create(&model.Video{AuthorID: 1, Title: "Video 1", PlayURL: "url1", CoverURL: "cover1"})
	db.Create(&model.Video{AuthorID: 2, Title: "Video 2", PlayURL: "url2", CoverURL: "cover2"})

	// Insert Comment fake data
	db.Create(&model.Comment{UserID: 1, VideoID: 1, Content: "Great video!", ActionType: "1"})
	db.Create(&model.Comment{UserID: 2, VideoID: 2, Content: "Nice work!", ActionType: "1"})

	// Insert Like fake data
	db.Create(&model.Like{UserID: 1, VideoID: 1, Liked: 1})
	db.Create(&model.Like{UserID: 2, VideoID: 2, Liked: 1})

	// Insert Message fake data
	db.Create(&model.Message{SenderID: 1, ReceiverID: 2, Content: "Hello!", ActionType: "1"})
	db.Create(&model.Message{SenderID: 2, ReceiverID: 1, Content: "Hi!", ActionType: "1"})

	// Insert Relation fake data
	db.Create(&model.Relation{UserID: 1, FollowingID: 2, Followed: 1})
	db.Create(&model.Relation{UserID: 2, FollowingID: 1, Followed: 0})

	//	打印插入数据成功
	fmt.Println("Insert fake data success!")

}
