package main

import (
	"fmt"
	"log"

	"github.com/RaymondCode/simple-demo/dao"
)

type UserCounts struct {
	FollowingCount int
	FansCount      int
}

func GetUserCounts(userIDs []uint) (map[uint]*UserCounts, error) {
	// 获取关注数
	var followingCounts []*UserCount
	err := dao.Relation.Select(dao.Relation.UserID, dao.Relation.ID.Count().As("count")).
		Where(dao.Relation.UserID.In(userIDs...), dao.Relation.Followed.Eq(1), dao.Relation.DeletedAt.IsNull()).
		Group(dao.Relation.UserID).
		Scan(&followingCounts)
	if err != nil {
		return nil, fmt.Errorf("无法获取用户的关注数量: %v", err)
	}

	// 获取粉丝数
	var fansCounts []*UserCount
	err = dao.Relation.Select(dao.Relation.FollowingID.As("user_id"), dao.Relation.ID.Count().As("count")).
		Where(dao.Relation.FollowingID.In(userIDs...), dao.Relation.Followed.Eq(1), dao.Relation.DeletedAt.IsNull()).
		Group(dao.Relation.FollowingID).
		Scan(&fansCounts)
	if err != nil {
		return nil, fmt.Errorf("无法获取用户的粉丝数量: %v", err)
	}

	// 创建一个映射以快速查找结果
	resultsMap := make(map[uint]*UserCounts)
	for _, count := range followingCounts {
		resultsMap[count.UserID] = &UserCounts{FollowingCount: count.Count}
	}
	for _, count := range fansCounts {
		if _, exists := resultsMap[count.UserID]; exists {
			resultsMap[count.UserID].FansCount = count.Count
		} else {
			resultsMap[count.UserID] = &UserCounts{FansCount: count.Count}
		}
	}

	// 确保所有传入的用户ID都包含在结果中
	for _, userID := range userIDs {
		if _, exists := resultsMap[userID]; !exists {
			resultsMap[userID] = &UserCounts{FollowingCount: 0, FansCount: 0}
		}
	}
	return resultsMap, nil
}

func main() {
	userCounts, err := GetUserCounts()
	if err != nil {
		log.Fatal(err)
	}

	// 遍历映射并打印每个用户的关注者和关注数
	for userID, userCount := range userCounts {
		fmt.Printf("User ID: %d\n", userID)
		fmt.Printf("Followers: %d\n", userCount.Followers)
		fmt.Printf("Following: %d\n", userCount.Following)
		fmt.Println("-----------------------------")
	}
}
