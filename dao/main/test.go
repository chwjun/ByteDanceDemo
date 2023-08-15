package main

import "github.com/RaymondCode/simple-demo/dao"

func IsVideoLikedByUser(userID uint, videoID uint) (bool, error) {
	likes := dao.Like

	count, err := likes.Where(likes.UserID.Eq(userID), likes.VideoID.Eq(videoID), likes.Liked.Eq(1)).Count()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func main() {
	println(IsVideoLikedByUser(1, 1))
}
