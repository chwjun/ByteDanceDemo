// Package casbin @Author: youngalone [2023/8/9]
package casbin

import (
	"github.com/RaymondCode/simple-demo/database"
	"github.com/casbin/casbin/v2"
	gormAdapter "github.com/casbin/gorm-adapter/v2"
	"github.com/gookit/slog"
)

func Setup() {
	e, err := GetCasbin()
	if err != nil {
		slog.Fatalf("casbin加载失败 %v", err)
	}
	{
		_, _ = e.AddPolicy("tourist", "/douyin/feed/", "GET")
		_, _ = e.AddPolicy("tourist", "/douyin/user/register/", "POST")
		_, _ = e.AddPolicy("tourist", "/douyin/user/login/", "POST")
	}

	{
		_, _ = e.AddPolicy("common_user", "/douyin/feed/", "GET")
		_, _ = e.AddPolicy("common_user", "/douyin/user/", "GET")
		_, _ = e.AddPolicy("common_user", "/publish/action/", "POST")
		_, _ = e.AddPolicy("common_user", "/publish/list/", "GET")
		_, _ = e.AddPolicy("common_user", "/favorite/action/", "POST")
		_, _ = e.AddPolicy("common_user", "/favorite/list/", "GET")
		_, _ = e.AddPolicy("common_user", "/comment/action/", "POST")
		_, _ = e.AddPolicy("common_user", "/relation/action/", "POST")
		_, _ = e.AddPolicy("common_user", "/relation/follow/list/", "GET")
		_, _ = e.AddPolicy("common_user", "/message/chat/", "GET")
		_, _ = e.AddPolicy("common_user", "/message/action/", "POST")
	}

	{
		_, _ = e.AddPolicy("admin", "/douyin/feed/", "GET")
		_, _ = e.AddPolicy("admin", "/douyin/user/register/", "POST")
		_, _ = e.AddPolicy("admin", "/douyin/user/login/", "POST")
		_, _ = e.AddPolicy("admin", "/douyin/user/", "GET")
		_, _ = e.AddPolicy("admin", "/publish/action/", "POST")
		_, _ = e.AddPolicy("admin", "/publish/list/", "GET")
		_, _ = e.AddPolicy("admin", "/favorite/action/", "POST")
		_, _ = e.AddPolicy("admin", "/favorite/list/", "GET")
		_, _ = e.AddPolicy("admin", "/comment/action/", "POST")
		_, _ = e.AddPolicy("admin", "/comment/list/", "GET")
		_, _ = e.AddPolicy("admin", "/relation/action/", "POST")
		_, _ = e.AddPolicy("admin", "/relation/follow/list/", "GET")
		_, _ = e.AddPolicy("admin", "/relation/follower/list/", "GET")
		_, _ = e.AddPolicy("admin", "/relation/friend/list/", "GET")
		_, _ = e.AddPolicy("admin", "/message/chat/", "GET")
		_, _ = e.AddPolicy("common_user", "/message/action/", "POST")
	}
}

func GetCasbin() (*casbin.Enforcer, error) {
	adapter, err := gormAdapter.NewAdapter("mysql", database.DSN, true)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		return nil, err
	}
	if err := e.LoadPolicy(); err == nil {
		return e, err
	} else {
		slog.Fatalf("权限管理系统加载失败 %v", err)
		return nil, err
	}
}
