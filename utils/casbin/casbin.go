// Package casbin @Author: youngalone [2023/8/9]
package casbin

import (
	"bytedancedemo/database/mysql"
	"github.com/casbin/casbin/v2"
	gormAdapter "github.com/casbin/gorm-adapter/v2"
	"go.uber.org/zap"
)

func Setup() {
	e, err := GetCasbin()
	if err != nil {
		zap.L().Error("权限管理系统加载失败", zap.Error(err))
	}
	{
		_, _ = e.AddPolicy("tourist", "/douyin/feed/", "GET")
		_, _ = e.AddPolicy("tourist", "/douyin/user/register/", "POST")
		_, _ = e.AddPolicy("tourist", "/douyin/user/login/", "POST")
	}

	{
		_, _ = e.AddPolicy("common_user", "/douyin/feed/", "GET")
		_, _ = e.AddPolicy("common_user", "/douyin/user/", "GET")
		_, _ = e.AddPolicy("common_user", "/douyin/publish/action/", "POST")
		_, _ = e.AddPolicy("common_user", "/douyin/publish/list/", "GET")
		_, _ = e.AddPolicy("common_user", "/douyin/favorite/action/", "POST")
		_, _ = e.AddPolicy("common_user", "/douyin/favorite/list/", "GET")
		_, _ = e.AddPolicy("common_user", "/douyin/comment/action/", "POST")
		_, _ = e.AddPolicy("common_user", "/douyin/relation/action/", "POST")
		_, _ = e.AddPolicy("common_user", "/douyin/relation/follow/list/", "GET")
		_, _ = e.AddPolicy("common_user", "/douyin/message/chat/", "GET")
		_, _ = e.AddPolicy("common_user", "/douyin/message/action/", "POST")
	}

	{
		_, _ = e.AddPolicy("admin", "/douyin/feed/", "GET")
		_, _ = e.AddPolicy("admin", "/douyin/user/register/", "POST")
		_, _ = e.AddPolicy("admin", "/douyin/user/login/", "POST")
		_, _ = e.AddPolicy("admin", "/douyin/user/", "GET")
		_, _ = e.AddPolicy("admin", "/douyin/publish/action/", "POST")
		_, _ = e.AddPolicy("admin", "/douyin/publish/list/", "GET")
		_, _ = e.AddPolicy("admin", "/douyin/favorite/action/", "POST")
		_, _ = e.AddPolicy("admin", "/douyin/favorite/list/", "GET")
		_, _ = e.AddPolicy("admin", "/douyin/comment/action/", "POST")
		_, _ = e.AddPolicy("admin", "/douyin/comment/list/", "GET")
		_, _ = e.AddPolicy("admin", "/douyin/relation/action/", "POST")
		_, _ = e.AddPolicy("admin", "/douyin/relation/follow/list/", "GET")
		_, _ = e.AddPolicy("admin", "/douyin/relation/follower/list/", "GET")
		_, _ = e.AddPolicy("admin", "/douyin/relation/friend/list/", "GET")
		_, _ = e.AddPolicy("admin", "/douyin/message/chat/", "GET")
		_, _ = e.AddPolicy("admin", "/douyin/message/action/", "POST")
	}
}

func GetCasbin() (*casbin.Enforcer, error) {
	adapter, err := gormAdapter.NewAdapter("mysql", mysql.DSN, true)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		return nil, err
	}
	if err := e.LoadPolicy(); err == nil {
		return e, nil
	} else {
		return nil, err
	}
}
