package service

type FriendUser struct {
	User
	Avatar  string `json:"avatar"`
	Message string `json:"message,omitempty"`
	MsgType int64  `json:"msg_type"`
}

type FollowService interface {

	/*
		模块业务接口
	*/
	// 当前用户关注目标用户
	FollowAction(userId int64, targetId int64) (bool, error)
	// 当前用户取消对目标用户的关注
	CancelFollowAction(userId int64, targetId int64) (bool, error)
	// 获取当前用户的关注列表
	GetFollowings(userId int64) ([]User, error)
	// 获取当前用户的粉丝列表
	GetFollowers(userId int64) ([]User, error)
	// 获取好友
	GetFriends(userId int64) ([]FriendUser, error)

	/*
		模块对外服务接口
	*/
	// 根据用户id查询关注数
	GetFollowingCnt(userId int64) (int64, error)
	// 根据用户id查询粉丝数
	GetFollowerCnt(userId int64) (int64, error)
	// 判断当前登录用户是否关注了目标用户
	CheckIsFollowing(userId int64, targetId int64) (bool, error)
}
