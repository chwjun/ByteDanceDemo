package main

import (
	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	// Initialize the generator with configuration
	g := gen.NewGenerator(gen.Config{
		OutPath:       "dao", // output directory, default value is ./query
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		FieldNullable: true,
	})

	// MySQL 连接字符串
	dsn := "sample_douyin:sample_douyin@tcp(43.140.203.85:3306)/sample_douyin?charset=utf8&parseTime=True&loc=Local&timeout=10000ms"
	// Initialize a *gorm.DB instance with MySQL
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Use the above `*gorm.DB` instance to initialize the generator
	g.UseDB(db)

	// Generate default DAO interface for those specified structs
	g.ApplyBasic(
		model.User{},
		model.Video{},
		model.Comment{},
		model.Like{},
		model.Message{},
		model.Relation{},
	)

	// Execute the generator
	g.Execute()
}
