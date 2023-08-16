// Package service @Author: youngalone [2023/8/8]
package service

import (
	"bytedancedemo/dao"
	"bytedancedemo/model"
	"github.com/gookit/slog"
	"github.com/spf13/viper"
	"sync"
)

type UserServiceImpl struct {
	// 这里需要关注模块 点赞模块 视频模块的配合
}

var (
	userServiceImpl *UserServiceImpl
	once            sync.Once
)

func GetUserServiceInstance() *UserServiceImpl {
	once.Do(func() {
		userServiceImpl = &UserServiceImpl{}
	})
	return userServiceImpl
}

func (usi *UserServiceImpl) InsertUser(user *model.User) (res *model.User, isSuccess bool) {
	u := dao.User
	err := u.Create(user)
	if err != nil {
		slog.Fatalf("新增用户失败 %v", err)
		return nil, false
	}
	resList, _ := u.Where(u.Name.Eq(user.Name), u.Password.Eq(user.Password)).Find()
	return resList[0], true
}

func (usi *UserServiceImpl) GetUserBasicByPassword(username string, password string) (user *model.User, isExist bool) {
	u := dao.User
	resList, err := u.Where(u.Name.Eq(username), u.Password.Eq(password)).Find()
	if err != nil {
		slog.Fatalf("查询用户失败 %v", err)
		return nil, false
	}
	if len(resList) == 0 {
		return nil, false
	}
	return resList[0], true
}

func (usi *UserServiceImpl) GetUserDetailsById(id int64, curID *int64) (*User, error) {
	user := &User{
		Id:              id,
		Name:            "抖音用户",
		Avatar:          viper.GetString("settings.oss.avatar"),
		BackgroundImage: viper.GetString("settings.oss.backgroundImage"),
		Signature:       viper.GetString("settings.oss.signature"),
	}
	u := dao.User
	resList, err := u.Where(u.ID.Eq(id)).Find()
	if err != nil {
		slog.Fatalf("查询用户失败 %v", err)
		return nil, err
	}
	user.Name = resList[0].Name
	// TODO 需要关注模块 点赞模块 视频模块的配合 获取剩余数据
	//var wg sync.WaitGroup
	//wg.Add(5)
	//if curID != nil {
	//	wg.Add(1)
	//}
	return user, nil
}
