package repository

import (
	"log"
	"sync"

	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/model"
)

type Follow struct {
	// Id          int64
	// UserId      int64
	// FollowingId int64
	// Followed    int8
	// CreatedAt   string
	// UpdatedAt   string
	model.Relation
}

func (Follow) TableName() string {
	return "relation"
}

type FollowDao struct {
}

var (
	followDao  *FollowDao
	followOnce sync.Once
)

// NewFollowDaoInstance 生成并返回followDao的单例对象。
func NewFollowDaoInstance() *FollowDao {
	followOnce.Do(
		func() {
			followDao = &FollowDao{}
		})
	return followDao
}

// // 单例模式
// func init() {
// 	followOnce.Do(func() {
// 		followDao = &FollowDao{}
// 	})
// }

// FindEverFollowing 给定当前用户和目标用户id，查看曾经是否有关注关系。
func (*FollowDao) FindEverFollowing(userId int64, targetId int64) (*model.Relation, error) {
	f := dao.Relation
	followList, err := f.Where(
		f.UserID.Eq(userId),
		f.FollowingID.Eq(targetId),
		f.Followed.In(0, 1),
	).Find()

	// 当查询出现错误时，日志打印 err msg，并 return err.
	if nil != err {
		// 当没查到记录报错时，不当做错误处理。
		if "record not found" == err.Error() {
			return nil, nil
		}
		log.Println(err.Error())
		return nil, err
	}

	// 如果查询到了关注关系，返回第一个关注关系
	if len(followList) > 0 {
		return followList[0], nil
	}

	return nil, nil

}

// InsertFollowRelation 给定用户和目标对象id，插入其关注关系。
func (*FollowDao) InsertFollowRelation(userId int64, targetId int64) (bool, error) {

	follow := &model.Relation{
		UserID:      userId,
		FollowingID: targetId,
		Followed:    1,
		//CreatedAt:   time.Now().Format(config.GO_STARTER_TIME),
	}

	// 将关注关系插入到数据库
	err := dao.Relation.Create(follow)
	if err != nil {
		log.Println(err) // 直接使用 err
		return false, err
	}

	return true, nil
}

// UpdateFollowRelation 给定用户和目标用户的id，更新他们的关系为取消关注或再次关注。
func (*FollowDao) UpdateFollowRelation(userId int64, targetId int64, followed int8) (bool, error) {
	f := dao.Relation

	// 构建更新的关注关系字段
	updates := map[string]interface{}{
		"Followed": followed,
	}

	// 更新关注关系状态
	_, result := f.Where(f.UserID.Eq(userId), f.FollowingID.Eq(targetId)).Updates(updates)
	if result != nil {
		log.Println(result)
		return false, result
	}

	return true, nil

}

// FindFollowRelation 给定当前用户和目标用户id，查询relation表是否存在关注关系
func (*FollowDao) FindFollowRelation(userId int64, targetId int64) (bool, error) {
	f := dao.Relation

	// 查询是否存在关注关系

	count, err := f.
		Where(f.UserID.Eq(userId), f.FollowingID.Eq(targetId), f.Followed.Eq(1)).
		Count()
	if err != nil {
		log.Println(err)
		return false, err
	}

	// 如果存在关注关系，返回 true
	if count > 0 {
		return true, nil
	}

	return false, nil
}

// GetFollowingsInfo 返回当前用户正在关注的用户信息列表，包括当前用户正在关注的用户ID列表和正在关注的用户总数
func (*FollowDao) GetFollowingsInfo(userId int64) ([]int64, int64, error) {

	f := dao.Relation

	// 查询正在关注的用户ID列表
	var followingIds []int64
	err := f.
		Where(f.UserID.Eq(userId), f.Followed.Eq(1)).
		Pluck(f.FollowingID, &followingIds)
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}

	// 查询正在关注的用户总数
	var totalFollowings int64
	totalFollowings, err = f.
		Where(f.UserID.Eq(userId), f.Followed.Eq(1)).
		Count()
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}

	return followingIds, totalFollowings, nil

}

// GetFollowersInfo 返回当前用户的粉丝用户信息列表，包括当前用户的粉丝用户ID列表和粉丝总数
func (*FollowDao) GetFollowersInfo(userId int64) ([]int64, int64, error) {

	f := dao.Relation

	// 查询粉丝用户ID列表
	var followerIds []int64
	err := f.
		Where(f.FollowingID.Eq(userId), f.Followed.Eq(1)).
		Pluck(f.UserID, &followerIds)
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}

	// 查询粉丝总数
	var totalFollowers int64
	totalFollowers, err = f.
		Where(f.FollowingID.Eq(userId), f.Followed.Eq(1)).
		Count()
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}

	return followerIds, totalFollowers, nil
}

func (*FollowDao) GetFriendsInfo(userId int64) ([]int64, int64, error) {

	friendId, friendCnt, err := followDao.GetFollowingsInfo(userId)

	if nil != err {
		log.Println(err.Error())
		return nil, -1, err
	}

	for i := 0; int64(i) < friendCnt; i++ {
		// 判断每一个登陆用户的关注用户是否关注了登陆用户，没关注就从集合里面剔除
		if flag, err1 := followDao.FindFollowRelation(friendId[i], userId); !flag {
			if err1 != nil {
				return nil, -1, err1
			}
			friendId = append(friendId[:i], friendId[i+1:]...)
			friendCnt--
			i--
		}

	}
	return friendId, friendCnt, nil

}

// GetFollowingCnt 给定当前用户id，查询relation表中该用户关注了多少人。
func (*FollowDao) GetFollowingCnt(userId int64) (int64, error) {
	f := dao.Relation

	// 查询正在关注的用户总数
	var followingCount int64
	followingCount, err := f.
		Where(f.UserID.Eq(userId), f.Followed.Eq(1)).
		Count()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return followingCount, nil
}

// GetFollowerCnt 给定当前用户id，查询relation表中该用户的粉丝数。
func (*FollowDao) GetFollowerCnt(userId int64) (int64, error) {
	f := dao.Relation

	// 查询粉丝总数
	var followerCount int64
	followerCount, err := f.
		Where(f.FollowingID.Eq(userId), f.Followed.Eq(1)).
		Count()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return followerCount, nil
}

// GetUserName 在user表中根据id查询用户姓名，放在followDao文件中并不妥当，后续可能废弃
func (*FollowDao) GetUserName(userId int64) (string, error) {
	u := dao.User

	var userName string
	err := u.
		Where(u.ID.Eq(userId)).
		Pluck(u.Name, &userName)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return userName, nil
}
