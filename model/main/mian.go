package main

import (
	"fmt"
	"log"

	"github.com/RaymondCode/simple-demo/model"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*创建表并插入假数据*/
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

	models := []interface{}{
		&model.User{},
		&model.Video{},
		&model.Comment{},
		&model.Like{},
		&model.Message{},
		&model.Relation{},
	}

	for _, model := range models {
		if db.Migrator().HasTable(model) {
			/*插入假数据*/
			insertFakeData(db)
			log.Fatalf("Table %T already exists", model)
		} else {
			err = db.AutoMigrate(model)
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("AutoMigrate success for %T!\n", model)
			}
		}
	}
	insertFakeData(db)
}

func insertFakeData(db *gorm.DB) {
	// Insert fake data into each table

	// Insert User fake data
	db.Create(&model.User{Name: "Alice", Password: "password123"})
	db.Create(&model.User{Name: "Bob", Password: "password456"})
	db.Create(&model.User{Name: "Charlie", Password: "password789"})
	db.Create(&model.User{Name: "David", Password: "password012"})
	db.Create(&model.User{Name: "Eve", Password: "password345"})
	db.Create(&model.User{Name: "Frank", Password: "password678"})
	db.Create(&model.User{Name: "Grace", Password: "password901"})
	db.Create(&model.User{Name: "Hannah", Password: "password234"})
	db.Create(&model.User{Name: "Isabel", Password: "password567"})
	db.Create(&model.User{Name: "Jack", Password: "password890"})

	// Insert Video fake data
	db.Create(&model.Video{AuthorID: 1, Title: "Video 1", PlayURL: "url1", CoverURL: "cover1"})
	db.Create(&model.Video{AuthorID: 2, Title: "Video 2", PlayURL: "url2", CoverURL: "cover2"})
	db.Create(&model.Video{AuthorID: 3, Title: "Video 3", PlayURL: "url3", CoverURL: "cover3"})
	db.Create(&model.Video{AuthorID: 4, Title: "Video 4", PlayURL: "url4", CoverURL: "cover4"})
	db.Create(&model.Video{AuthorID: 5, Title: "Video 5", PlayURL: "url5", CoverURL: "cover5"})
	db.Create(&model.Video{AuthorID: 6, Title: "Video 6", PlayURL: "url6", CoverURL: "cover6"})
	db.Create(&model.Video{AuthorID: 7, Title: "Video 7", PlayURL: "url7", CoverURL: "cover7"})
	db.Create(&model.Video{AuthorID: 8, Title: "Video 8", PlayURL: "url8", CoverURL: "cover8"})
	db.Create(&model.Video{AuthorID: 9, Title: "Video 9", PlayURL: "url9", CoverURL: "cover9"})
	db.Create(&model.Video{AuthorID: 10, Title: "Video 10", PlayURL: "url10", CoverURL: "cover10"})

	// Insert Comment fake data
	db.Create(&model.Comment{UserID: 1, VideoID: 1, Content: "Great video!", ActionType: 1})
	db.Create(&model.Comment{UserID: 2, VideoID: 2, Content: "Nice work!", ActionType: 1})
	db.Create(&model.Comment{UserID: 3, VideoID: 3, Content: "Fantastic!", ActionType: 1})
	db.Create(&model.Comment{UserID: 4, VideoID: 4, Content: "Well done!", ActionType: 1})
	db.Create(&model.Comment{UserID: 5, VideoID: 5, Content: "Loved it!", ActionType: 1})
	db.Create(&model.Comment{UserID: 6, VideoID: 6, Content: "Keep it up!", ActionType: 1})
	db.Create(&model.Comment{UserID: 7, VideoID: 7, Content: "Amazing work!", ActionType: 1})
	db.Create(&model.Comment{UserID: 1, VideoID: 8, Content: "So creative!", ActionType: 1})
	db.Create(&model.Comment{UserID: 2, VideoID: 9, Content: "Inspiring!", ActionType: 1})
	db.Create(&model.Comment{UserID: 9, VideoID: 10, Content: "Brilliant!", ActionType: 1})

	// Insert Like fake data
	db.Create(&model.Like{UserID: 1, VideoID: 1, Liked: 1})
	db.Create(&model.Like{UserID: 2, VideoID: 2, Liked: 1})
	db.Create(&model.Like{UserID: 3, VideoID: 1, Liked: 1})
	db.Create(&model.Like{UserID: 4, VideoID: 2, Liked: 1})
	db.Create(&model.Like{UserID: 5, VideoID: 3, Liked: 0})
	db.Create(&model.Like{UserID: 6, VideoID: 4, Liked: 1})
	db.Create(&model.Like{UserID: 7, VideoID: 5, Liked: 0})
	db.Create(&model.Like{UserID: 8, VideoID: 6, Liked: 1})
	db.Create(&model.Like{UserID: 9, VideoID: 7, Liked: 1})
	db.Create(&model.Like{UserID: 10, VideoID: 8, Liked: 0})

	// Insert Message fake data
	db.Create(&model.Message{SenderID: 1, ReceiverID: 2, Content: "Hello!", ActionType: 1})
	db.Create(&model.Message{SenderID: 2, ReceiverID: 1, Content: "Hi!", ActionType: 1})
	db.Create(&model.Message{SenderID: 3, ReceiverID: 4, Content: "How are you?", ActionType: 1})
	db.Create(&model.Message{SenderID: 4, ReceiverID: 3, Content: "I'm fine, thanks!", ActionType: 1})
	db.Create(&model.Message{SenderID: 5, ReceiverID: 6, Content: "Good morning!", ActionType: 1})
	db.Create(&model.Message{SenderID: 6, ReceiverID: 5, Content: "Good evening!", ActionType: 2}) // This message is retracted
	db.Create(&model.Message{SenderID: 7, ReceiverID: 8, Content: "Happy birthday!", ActionType: 1})
	db.Create(&model.Message{SenderID: 8, ReceiverID: 7, Content: "Thank you!", ActionType: 1})
	db.Create(&model.Message{SenderID: 9, ReceiverID: 10, Content: "See you soon!", ActionType: 1})
	db.Create(&model.Message{SenderID: 10, ReceiverID: 9, Content: "Looking forward!", ActionType: 1})

	// Insert Relation fake data
	db.Create(&model.Relation{UserID: 1, FollowingID: 2, Followed: 1})
	db.Create(&model.Relation{UserID: 2, FollowingID: 1, Followed: 0}) // Insert additional Relation fake data
	db.Create(&model.Relation{UserID: 3, FollowingID: 4, Followed: 1})
	db.Create(&model.Relation{UserID: 4, FollowingID: 3, Followed: 1})
	db.Create(&model.Relation{UserID: 5, FollowingID: 6, Followed: 0})
	db.Create(&model.Relation{UserID: 6, FollowingID: 5, Followed: 1})
	db.Create(&model.Relation{UserID: 7, FollowingID: 8, Followed: 1})
	db.Create(&model.Relation{UserID: 8, FollowingID: 9, Followed: 0})
	db.Create(&model.Relation{UserID: 9, FollowingID: 10, Followed: 1})
	db.Create(&model.Relation{UserID: 10, FollowingID: 1, Followed: 1})

	//	打印插入数据成功
	fmt.Println("Insert fake data success!")

}
