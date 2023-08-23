package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"log"
)

// Package gen @Author: youngalone [2023/8/7]

func main() {
	dsn := "sample_douyin:sample_douyin@tcp(43.140.203.85:3306)/sample_douyin?charset=utf8&parseTime=True&loc=Local&timeout=10000ms"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	g := gen.NewGenerator(gen.Config{
		OutPath:           "./dao",
		ModelPkgPath:      "./model",
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		FieldNullable:     false,
		FieldCoverable:    false,
		FieldSignable:     false,
		FieldWithIndexTag: false,
		FieldWithTypeTag:  true,
	})
	g.UseDB(db)
	allModel := g.GenerateAllTable()
	g.ApplyBasic(allModel...)
	g.Execute()
}
