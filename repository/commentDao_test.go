// Package repository @Author: youngalone [2023/8/16]
package repository

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/database"
	"github.com/RaymondCode/simple-demo/model"
	"testing"
)

func TestDeleteComment(t *testing.T) {
	config.Init("../config/settings.yml")
	database.Init()
	dao.SetDefault(database.DB)
	err := DeleteComment(25)
	fmt.Printf("err = %v\n", err)
}

func TestGetCommentCnt(t *testing.T) {
	config.Init("../config/settings.yml")
	database.Init()
	dao.SetDefault(database.DB)
	res, err := GetCommentCnt(1)
	fmt.Printf("res = %v\n", res)
	fmt.Printf("err = %v\n", err)
}

func TestGetCommentList(t *testing.T) {
	config.Init("../config/settings.yml")
	database.Init()
	dao.SetDefault(database.DB)
	list, err := GetCommentList(1)
	for i, comment := range list {
		fmt.Printf("comment %d = %v\n", i, comment)
	}
	fmt.Printf("err = %v\n", err)
}

func TestInsertComment(t *testing.T) {
	config.Init("../config/settings.yml")
	database.Init()
	dao.SetDefault(database.DB)
	comment, err := InsertComment(model.Comment{
		UserID:     3,
		VideoID:    3,
		Content:    "测试插入评论",
		ActionType: 1,
	})
	fmt.Printf("err = %v\n", err)
	fmt.Printf("res = %v\n", comment)
}
