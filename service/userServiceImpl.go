// Package service @Author: youngalone [2023/8/8]
package service

import (
	"bytedancedemo/dao"
	"bytedancedemo/model"
	"bytedancedemo/utils"
	"errors"
	"go.uber.org/zap"
	"sync"
)

type UserServiceImpl struct {
	// 这里需要关注模块 点赞模块 视频模块的配合
	VideoService
	FollowService
}

var (
	userServiceImpl *UserServiceImpl
	once            sync.Once
)

func GetUserServiceInstance() *UserServiceImpl {
	once.Do(func() {
		userServiceImpl = &UserServiceImpl{
			VideoService:  &VideoServiceImp{},
			FollowService: &FollowServiceImp{},
		}
	})
	return userServiceImpl
}

func (usi *UserServiceImpl) InsertUser(user *model.User) (res *model.User, isSuccess bool) {
	u := dao.User
	err := u.Create(user)
	if err != nil {

		zap.L().Fatal("新增用户失败 ", zap.String("err", err.Error()))
		return nil, false
	}
	resList, _ := u.Where(u.Name.Eq(user.Name), u.Password.Eq(user.Password)).Find()
	return resList[0], true
}

func (usi *UserServiceImpl) GetUserBasicByPassword(username string, password string) (user *model.User, isExist bool) {
	u := dao.User
	resList, err := u.Where(u.Name.Eq(username), u.Password.Eq(password)).Find()
	if err != nil {
		zap.L().Fatal("查询用户失败", zap.String("err", err.Error()))
		return nil, false
	}
	if len(resList) == 0 {
		zap.L().Warn("未查询到用户", zap.Error(errors.New("用户名或密码错误")))
		return nil, false
	}
	return resList[0], true
}

func (usi *UserServiceImpl) GetUserDetailsById(id int64, curID *int64) (*User, error) {
	user := &User{
		Id:   id,
		Name: "抖音用户",
	}
	u := dao.User
	resList, err := u.Where(u.ID.Eq(id)).Find()

	if err != nil || len(resList) == 0 {
		zap.L().Fatal("查询用户失败 ", zap.String("err", err.Error()))
		return nil, err
	}
	user.Name = resList[0].Name
	user.Avatar = resList[0].Avatar
	user.BackgroundImage = resList[0].BackgroundImage
	user.Signature = resList[0].Signature

	userService := GetUserServiceInstance()

	var wg sync.WaitGroup
	wg.Add(5)
	if curID != nil {
		wg.Add(1)
		// 判断是否关注
		go func() {
			isFollow, err := userService.CheckIsFollowing(id, *curID)
			if err != nil {
				wg.Done()
				return
			}
			user.IsFollow = isFollow
			wg.Done()
		}()
	}

	// 获取作品数
	go func() {
		workCnt, err := userService.GetVideoCountByAuthorID(id)
		if err != nil {
			wg.Done()
			return
		}
		user.WorkCount = workCnt
		wg.Done()
	}()

	// 获取关注数
	go func() {
		cnt, err := userService.GetFollowingCnt(id)
		if err != nil {
			wg.Done()
			return
		}
		user.FollowCount = cnt
		wg.Done()
	}()

	// 获取粉丝数
	go func() {
		cnt, err := userService.GetFollowerCnt(id)
		if err != nil {
			wg.Done()
			return
		}
		user.FollowerCount = cnt
		wg.Done()
	}()

	// 获取获赞数
	go func() {
		likes, err := GetUserTotalReceivedLikes([]int64{id})
		if err != nil {
			wg.Done()
			return
		}
		user.TotalFavorited = likes[id]
		wg.Done()
	}()

	// 获取点赞数
	go func() {
		favorites, err := utils.GetUserFavorites([]int64{id})
		if err != nil {
			wg.Done()
			return
		}
		user.FavoriteCount = favorites[id]
		wg.Done()
	}()

	wg.Wait()
	return user, nil
}

// GetUserName 在user表中根据id查询用户姓名
func (usi *UserServiceImpl) GetUserName(userId int64) (string, error) {
	u := dao.User

	var userName string
	err := u.
		Where(u.ID.Eq(userId)).
		Pluck(u.Name, &userName)
	if err != nil {
		zap.L().Error("查询用户名出错", zap.Error(err))
		return "", err
	}

	return userName, nil
}
