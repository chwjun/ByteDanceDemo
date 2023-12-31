package repository

import (
	"bytedancedemo/config"
	"bytedancedemo/dao"
	"bytedancedemo/database/mysql"
	"fmt"
	"log"
	"testing"
)

func TestFollowDao_FindEverFollowing(t *testing.T) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)
	follow, err := followDao.FindEverFollowing(1, 2)
	if err == nil {
		log.Default()
	}
	fmt.Print(follow)
}

func TestFollowDao_InsertFollowRelation(t *testing.T) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)
	isbool, err := followDao.InsertFollowRelation(2, 3)
	if err == nil {
		log.Default()
	}
	fmt.Print(isbool)
}

func TestFollowDao_UpdateFollowRelation(t *testing.T) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)
	isbool, err := followDao.UpdateFollowRelation(2, 3, 0)
	if err == nil {
		log.Default()
	}
	fmt.Print(isbool)

}

func TestFollowDao_FindFollowRelation(t *testing.T) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)
	isbool, err := followDao.FindFollowRelation(2, 3)
	if err == nil {
		log.Default()
	}
	fmt.Print(isbool)

}

func TestFollowDao_GetFollowingsInfo(t *testing.T) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)
	followingsId, followingsCnt, err := followDao.GetFollowingsInfo(1)

	if err != nil {
		log.Default()
	}

	fmt.Println(followingsId)
	fmt.Println(followingsCnt)

}

func TestFollowDao_GetFollowersInfo(t *testing.T) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)
	followersId, followersCnt, err := followDao.GetFollowersInfo(1)

	if err != nil {
		log.Default()
	}

	fmt.Println(followersId)
	fmt.Println(followersCnt)

}

func TestFollowDao_GetFriendsInfo(t *testing.T) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)
	friendId, friendCnt, err := followDao.GetFriendsInfo(1)
	if err != nil {
		log.Default()
	}
	fmt.Println(friendId)
	fmt.Println(friendCnt)

}

func TestFollowDao_GetFollowingCnt(t *testing.T) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)
	followingCount, err := followDao.GetFollowingCnt(1)
	if err != nil {
		log.Default()
	}
	fmt.Println(followingCount)

}

func TestFollowDao_GetFollowerCnt(t *testing.T) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)
	followerCount, err := followDao.GetFollowerCnt(1)
	if err != nil {
		log.Default()
	}
	fmt.Println(followerCount)

}
