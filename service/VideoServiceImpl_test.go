package service

import (
	"bytedancedemo/config"
	"bytedancedemo/dao"
	"bytedancedemo/database"
	"testing"
	"time"
	// "gorm.io/driver/mysql"
)

// 吃的是单例模式
func TestNewVSIInstance(t *testing.T) {
	vsi1 := NewVSIInstance()
	vsi2 := NewVSIInstance()
	if vsi1 != vsi2 {
		t.Error("单例测试出错")
	} else {
		t.Logf("单例测试成功")
	}
}

// go test -v .\service\VideoServiceImpl_test.go .\service\VideoServiceImpl.go .\service\VideoService.go .\service\UserService.go .\service\UserServiceImpl.go
// 测试Feed接口
func TestVideoServiceImpl_Feed(t *testing.T) {
	config.Init("../config/settings.yml")
	database.Init()
	dao.SetDefault(database.DB)
	time1 := time.Date(2023, 8, 19, 18, 20, 0, 0, time.Now().Location())
	//timeNow := time.Now()
	// 正确的4种情况，错误的有1种（userid不存在）
	type test struct {
		name          string
		latest_time   time.Time
		user_id       int64
		want_VideoID  []int64
		want_authorID []int64
	}
	tests := make([]test, 5)
	tests[1].name = "测试1，有latesttime，无userid"
	tests[1].user_id = 0
	tests[1].latest_time = time1
	tests[1].want_VideoID = []int64{3, 4, 5, 6, 7, 8, 9, 10}
	tests[1].want_authorID = []int64{3, 4, 5, 6, 7, 8, 9, 10}
	t.Run(tests[1].name, func(t *testing.T) {
		vsi := &VideoServiceImp{}
		videolist, _, err := vsi.Feed(tests[1].latest_time, int(tests[1].user_id))
		if err != nil {
			t.Errorf(err.Error())
		} else if len(videolist) == 0 {
			t.Errorf("没有找到视频列表")
		} else if videolist[0].Id != int(tests[1].want_VideoID[7]) {
			t.Errorf("视频编号错误")
		} else {
			t.Logf("成功！")
		}
	})
}
