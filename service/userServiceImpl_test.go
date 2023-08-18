// Package service @Author: youngalone [2023/8/18]
package service

import (
	"bytedancedemo/config"
	"bytedancedemo/dao"
	"bytedancedemo/database/mysql"
	"bytedancedemo/model"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// 测试单例模式
func TestGetUserServiceInstance(t *testing.T) {
	usi1 := GetUserServiceInstance()
	usi2 := GetUserServiceInstance()
	if usi1 != usi2 {
		t.Errorf("单例模式出错")
	}
}

// 测试根据用户名和密码查询账户
func TestUserServiceImpl_GetUserBasicByPassword(t *testing.T) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name        string
		args        args
		wantUserID  int64
		wantIsExist bool
	}{
		{"测试一 合法账户测试", args{"admin", "123456"}, 0, true},
		{"测试二 非法账户测试", args{"admin", "1234567"}, 0, false},
	}
	t.Run(tests[0].name, func(t *testing.T) {
		usi := &UserServiceImpl{}
		gotUser, gotIsExist := usi.GetUserBasicByPassword(tests[0].args.username, tests[0].args.password)
		if gotIsExist != tests[0].wantIsExist {
			t.Errorf("未成功找到用户账号")
		}
		if gotUser.ID != tests[0].wantUserID {
			t.Errorf("找到错误的用户账号")
		}
	})
	t.Run(tests[1].name, func(t *testing.T) {
		usi := &UserServiceImpl{}
		gotUser, gotIsExist := usi.GetUserBasicByPassword(tests[1].args.username, tests[1].args.password)
		if gotIsExist != tests[1].wantIsExist || gotUser != nil {
			t.Errorf("找到错误的用户账号")
		}
	})
}

// 测试根据用户ID获取用户详情
func TestUserServiceImpl_GetUserDetailsById(t *testing.T) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)
	var curID1 int64 = 1
	type args struct {
		id    int64
		curID *int64
	}
	tests := []struct {
		name       string
		args       args
		wantName   string
		wantFollow bool
		wantErr    bool
	}{
		{"测试一 获取用户", args{1, nil}, "Alice", false, false},
		{"测试二 获取目标用户", args{1, &curID1}, "Alice", true, false},
	}
	t.Run(tests[0].name, func(t *testing.T) {
		usi := &UserServiceImpl{}
		got, gotErr := usi.GetUserDetailsById(tests[0].args.id, tests[0].args.curID)
		if (gotErr != nil) != tests[0].wantErr {
			t.Errorf("查询用户时出错")
		}
		if got.Name != tests[0].wantName {
			t.Errorf("查询到错误用户")
		}
		if got.IsFollow != tests[0].wantFollow {
			t.Errorf("关注状态出错")
		}
	})
	t.Run(tests[1].name, func(t *testing.T) {
		usi := &UserServiceImpl{}
		got, gotErr := usi.GetUserDetailsById(tests[1].args.id, tests[1].args.curID)
		if (gotErr != nil) != tests[1].wantErr {
			t.Errorf("查询用户时出错")
		}
		if got.Name != tests[1].wantName {
			t.Errorf("查询到错误用户")
		}
		if got.IsFollow != tests[1].wantFollow {
			t.Errorf("关注状态出错")
		}
	})
}

// 测试插入用户
func TestUserServiceImpl_InsertUser(t *testing.T) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)
	type args struct {
		user *model.User
	}
	rand.New(rand.NewSource(time.Now().Unix()))
	tests := []struct {
		name          string
		args          args
		wantIsSuccess bool
	}{
		{"测试一 正常使用测试", args{&model.User{Name: fmt.Sprintf("用户%d", rand.Intn(99_999_999_999)), Password: "123456", Role: "common_role"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usi := &UserServiceImpl{}
			gotUser, gotIsSuccess := usi.InsertUser(tt.args.user)
			if gotUser.Name != tt.args.user.Name {
				t.Errorf("InsertUser() gotUser.Name = %v, want %v", gotUser.Name, tt.args.user.Name)
			}
			if gotIsSuccess != tt.wantIsSuccess {
				t.Errorf("InsertUser() gotIsSuccess = %v, want %v", gotIsSuccess, tt.wantIsSuccess)
			}
		})
	}
}
