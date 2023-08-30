package service

import (
	"bytedancedemo/utils"
	"fmt"
	"testing"
)

func BenchmarkFavoriteList(b *testing.B) {
	// 创建服务实例
	service := &FavoriteServiceImpl{
		// 您可以在此处设置任何所需的模拟或实际依赖项
	}

	// 设置测试数据嗯，
	userID := int64(1)

	// 运行基准测试
	b.ResetTimer() // 重置计时器以排除设置过程的时间
	for i := 0; i < b.N; i++ {
		_, _ = service.FavoriteList(userID)
	}
}

func TestGetFavoriteVideoInfoByUserID(t *testing.T) {
	// 创建服务实例
	service := &FavoriteServiceImpl{}

	// 使用一个有效的用户ID进行测试
	testUserID := int64(1) // 请确保这个ID在你的测试环境中是有效的

	// 调用要测试的函数
	videos, err := service.GetFavoriteVideoInfoByUserID(testUserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 进行断言，验证结果是否符合预期
	// 注意：以下断言应根据你的实际需求和业务逻辑进行调整
	if len(videos) == 0 {
		t.Errorf("expected at least one video, got none")
	}

	// 其他可能的断言，例如验证返回的视频信息
	for _, video := range videos {
		if video.ID <= 0 {
			t.Errorf("expected positive video ID, got %d", video.ID)
		}
		if video.Author.Id <= 0 {
			t.Errorf("expected positive author ID, got %d", video.Author.Id)
		}
		// 可以继续添加其他检查
	}
}

func TestGetUserDetailsByIDs(t *testing.T) {
	// 这里可能需要进行数据库连接和初始化
	// 假设你已经有了一个已初始化的数据库连接 Db

	// 选择一些用户ID作为测试参数
	userIDs := []int64{1, 2, 3}

	// 调用函数
	userDetails, err := GetUserDetailsByIDs(userIDs)
	if err != nil {
		t.Errorf("出现错误: %v", err)
		return
	}

	// 打印结果
	for userID, detail := range userDetails {
		fmt.Printf("用户ID: %v, 用户名: %v, 头像: %v, 背景图片: %v, 个性签名: %v\n",
			userID, detail.Name, detail.Avatar, detail.BackgroundImage, detail.Signature)
	}
}
func TestGetUserInfoByIDs(t *testing.T) {
	// 创建一个新的 FavoriteServiceImpl
	service := &FavoriteServiceImpl{}

	// 使用测试用户ID和用户ID数组调用方法
	requestingUserID := int64(123)
	userIDs := []int64{1, 2}
	expectedUsersCount := len(userIDs)

	users, err := service.GetUserInfoByIDs(requestingUserID, userIDs)

	// 检查是否存在错误
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// 检查返回的用户数量
	if len(users) != expectedUsersCount {
		t.Fatalf("Expected %v users, got %v", expectedUsersCount, len(users))
	}

	// 打印每个用户的每个字段
	for i, user := range users {
		fmt.Printf("User %d:\n", i+1)
		fmt.Printf("  ID: %v\n", user.Id)
		fmt.Printf("  Name: %s\n", user.Name)
		fmt.Printf("  FollowCount: %v\n", user.FollowCount)
		fmt.Printf("  FollowerCount: %v\n", user.FollowerCount)
		fmt.Printf("  IsFollow: %v\n", user.IsFollow)
		fmt.Printf("  Avatar: %s\n", user.Avatar)
		fmt.Printf("  BackgroundImage: %s\n", user.BackgroundImage)
		fmt.Printf("  Signature: %s\n", user.Signature)
		fmt.Printf("  TotalFavorited: %v\n", user.TotalFavorited)
		fmt.Printf("  WorkCount: %v\n", user.WorkCount)
		fmt.Printf("  FavoriteCount: %v\n", user.FavoriteCount)
	}
}

func TestGetCommentCounts(t *testing.T) {
	// 创建一些测试数据
	videoIDs := []int64{1, 2, 3}

	// 调用测试的函数
	_, err := GetCommentCounts(videoIDs)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// 添加更多的检查，以验证 GetCommentCounts 的返回值
	// ...
}
func TestGetUserCounts(t *testing.T) {
	userIDs := []int64{1, 2}

	// 测试 "following" 类型的计数
	resultFollowing, err := GetUserCounts(userIDs, "following")
	if err != nil {
		t.Fatalf("Error getting following counts: %v", err)
	}
	followingMap := make(map[int64]int64)
	for _, count := range resultFollowing {
		followingMap[count.UserID] = count.Count
	}
	fmt.Println("Following counts for user IDs [1, 2]:")
	for _, userID := range userIDs {
		fmt.Printf("UserID: %v, Count: %v\n", userID, followingMap[userID])
	}

	// 测试 "fans" 类型的计数
	resultFans, err := GetUserCounts(userIDs, "fans")
	if err != nil {
		t.Fatalf("Error getting fans counts: %v", err)
	}
	fansMap := make(map[int64]int64)
	for _, count := range resultFans {
		fansMap[count.UserID] = count.Count
	}
	fmt.Println("Fans counts for user IDs [1, 2]:")
	for _, userID := range userIDs {
		fmt.Printf("UserID: %v, Count: %v\n", userID, fansMap[userID])
	}
}
func TestGetLikeCounts(t *testing.T) {
	// 输入的视频ID
	videoIDs := []int64{1, 2, 3}

	// 调用函数
	result, err := utils.GetUserFavorites(videoIDs)
	if err != nil {
		t.Fatalf("GetLikeCounts failed with error: %v", err)
	}

	// 检查结果
	if len(result) != len(videoIDs) {
		t.Fatalf("Expected result length %d, but got %d", len(videoIDs), len(result))
	}

	for _, id := range videoIDs {
		if count, ok := result[id]; ok {
			fmt.Printf("Video ID %d has %d likes\n", id, count)
		} else {
			t.Fatalf("Expected video ID %d in result, but not found", id)
		}
	}
}
