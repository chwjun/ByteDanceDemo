package service

import (
	"bytedancedemo/config"
	"bytedancedemo/dao"
	"bytedancedemo/database/mysql"
	"testing"
	"time"
	// "gorm.io/driver/mysql"
)

// 测的是单例模式
func TestNewVSIInstance(t *testing.T) {
	vsi1 := NewVSIInstance()
	vsi2 := NewVSIInstance()
	if vsi1 != vsi2 {
		t.Error("单例测试出错")
	} else {
		t.Logf("单例测试成功")
	}
}

// go test -v .\service\VideoServiceImpl_test.go .\service\VideoServiceImpl.go .\service\VideoService.go .\service\userService.go .\service\UserServiceImpl.go
// 测试Feed接口
func TestVideoServiceImpl_Feed(t *testing.T) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)
	time1 := time.Date(2023, 8, 19, 18, 07, 46, 0, time.Now().Location())
	timeNow := time.Now()
	// 正确的4种情况，错误的有1种（userid不存在）
	type test struct {
		name          string
		latest_time   time.Time
		user_id       int64
		want_VideoID  int64
		want_authorID int64
	}
	tests := make([]test, 6)
	tests[1].name = "测试1，有latesttime，无userid"
	tests[1].user_id = 0
	tests[1].latest_time = time1
	tests[1].want_VideoID = 7
	tests[1].want_authorID = 7
	t.Run(tests[1].name, func(t *testing.T) {
		vsi := &VideoServiceImp{}
		videolist, _, err := vsi.Feed(tests[1].latest_time.UnixMilli(), tests[1].user_id)
		if err != nil {
			t.Errorf(err.Error())
		} else if len(videolist) == 0 {
			t.Errorf("没有找到视频列表")
		} else if videolist[0].Id != tests[1].want_VideoID {
			t.Logf("len of video list %v video list id %v", len(videolist), videolist[0].Id)
			t.Errorf("视频编号错误")
		} else {
			t.Logf("成功！")
		}
	})
	tests[2].name = "测试2，无latesttime，使用当前时间"
	tests[2].latest_time = timeNow
	tests[2].user_id = 0
	tests[2].want_VideoID = 10
	tests[2].want_authorID = 10
	t.Run(tests[2].name, func(t *testing.T) {
		vsi := &VideoServiceImp{}
		videolist, _, err := vsi.Feed(tests[2].latest_time.UnixMilli(), tests[2].user_id)
		if err != nil {
			t.Errorf(err.Error())
		} else if len(videolist) == 0 {
			t.Errorf("没有找到视频列表")
		} else if videolist[0].Id != tests[2].want_VideoID {
			t.Logf("len of video list %v video list id %v", len(videolist), videolist[0].Id)
			t.Errorf("视频编号错误")
		} else {
			t.Logf("成功！")
		}
	})
	// userid没用到，假数据
}
func TestVideoServiceImpl_Pbulishlist(t *testing.T) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)
	// 正确的4种情况，错误的有1种（userid不存在）
	type test struct {
		name          string
		user_id       int64
		want_VideoID  int64
		want_authorID int64
	}
	tests := make([]test, 6)
	tests[1].name = "测试1 无userid"
	tests[1].user_id = 0
	// tests[1].latest_time = time1
	tests[1].want_VideoID = 7
	tests[1].want_authorID = 7
	t.Run(tests[1].name, func(t *testing.T) {
		vsi := &VideoServiceImp{}
		videolist, err := vsi.PublishList(tests[1].user_id)
		if err != nil {
			t.Errorf(err.Error())
		} else if len(videolist) == 0 {
			t.Logf("没有视频id，运行成功")
		} else if videolist[0].Id != tests[1].want_VideoID {
			t.Logf("len of video list %v video list id %v", len(videolist), videolist[0].Id)
			t.Errorf("视频编号错误")
		}
	})
	tests[2].name = "测试2 有userID，ID合法"
	tests[2].user_id = 2
	tests[2].want_VideoID = 2
	tests[2].want_authorID = 2
	t.Run(tests[2].name, func(t *testing.T) {
		vsi := &VideoServiceImp{}
		videolist, err := vsi.PublishList(tests[2].user_id)
		if err != nil {
			t.Errorf(err.Error())
		} else if len(videolist) == 0 {
			t.Errorf("没有找到视频列表")
		} else if videolist[0].Id != tests[2].want_VideoID {
			t.Logf("len of video list %v video list id %v", len(videolist), videolist[0].Id)
			t.Errorf("视频编号错误")
		} else {
			t.Logf("成功！")
		}
	})
	tests[2].name = "测试3 有userID，ID不合法"
	tests[2].user_id = 11
	tests[2].want_VideoID = 10
	tests[2].want_authorID = 10
	t.Run(tests[2].name, func(t *testing.T) {
		vsi := &VideoServiceImp{}
		videolist, err := vsi.PublishList(tests[2].user_id)
		if err != nil {
			t.Errorf(err.Error())
		} else if len(videolist) == 0 {
			t.Logf("没有找到视频列表，成功")
		} else if videolist[0].Id != tests[2].want_VideoID {
			t.Logf("len of video list %v video list id %v", len(videolist), videolist[0].Id)
			t.Errorf("视频编号错误")
		}
	})
}

func BenchmarkNewVSIInstance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewVSIInstance()
	}
}

// 基准测试，测试Feed
func BenchmarkVideoServiceImpl_Feed(b *testing.B) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)

	vsi := NewVSIInstance()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		vsi.Feed(time.Now().UnixMilli(), 0)
	}
}

func BenchmarkVideoServiceImpl_Pbulishlist(b *testing.B) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)

	vsi := NewVSIInstance()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vsi.PublishList(int64(i))
	}
}
