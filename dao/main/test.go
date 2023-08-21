package main

import (
	"fmt"

	"github.com/RaymondCode/simple-demo/dao"
)

func GetUserVideoIDs(userID uint) ([]uint, error) {
	// 存储查询结果的切片
	var videoIDs []uint

	// 执行查询，并将结果扫描到 videoIDs 切片
	err := dao.Video.Select(dao.Video.ID).Where(dao.Video.AuthorID.Eq(userID), dao.Video.DeletedAt.IsNull()).Scan(&videoIDs)
	if err != nil {
		return nil, fmt.Errorf("无法获取用户的全部视频ID: %v", err)
	}

	return videoIDs, nil
}

func main() {
	userID := uint(1) // 换成你想查询的用户ID
	videoIDs, err := GetUserVideoIDs(userID)
	if err != nil {
		fmt.Println("获取视频ID失败:", err)
		return
	}

	fmt.Println("获取到的视频ID:")
	for _, id := range videoIDs {
		fmt.Println(id)
	}
}
