package main

import (
	"fmt"

	"github.com/RaymondCode/simple-demo/dao"
)

func GetUserTotalReceivedLikes(userIDs []uint) (map[uint]int64, error) {
	likesCount := make(map[uint]int64)

	// 定义一个结构体来保存查询结果
	type LikeResult struct {
		UserID uint
		Count  int64
	}

	var results []LikeResult
	err := dao.Like.
		Where(dao.Like.UserID.In(userIDs...), dao.Like.Liked.Eq(1), dao.Like.DeletedAt.IsNull()).
		Group(dao.Like.UserID).
		Select(dao.Like.UserID.As("user_id"), dao.Like.ID.Count().As("count")).
		Scan(&results)

	if err != nil {
		return nil, fmt.Errorf("无法获取用户总接收的喜欢数量: %v", err)
	}

	for _, result := range results {
		likesCount[result.UserID] = result.Count
	}

	// 对于没有获取到喜欢数量的用户，将他们添加到map中，并设置值为0
	for _, id := range userIDs {
		if _, ok := likesCount[id]; !ok {
			likesCount[id] = 0
		}
	}

	return likesCount, nil
}

func main() {

	// 这里假设你有一个用户ID的列表
	userIDs := []uint{1, 2, 3}

	// 调用 GetUserTotalReceivedLikes 函数
	likesCount, err := GetUserTotalReceivedLikes(userIDs)

	// 检查错误
	if err != nil {
		fmt.Println("获取用户喜欢数量失败:", err)
		return
	}

	// 打印结果
	fmt.Println("用户喜欢数量:")
	for userID, count := range likesCount {
		fmt.Printf("用户ID: %d, 喜欢数量: %d\n", userID, count)
	}

}
