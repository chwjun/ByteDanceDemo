package service

type Video struct {
	ID            int64  `json:"id"`
	Author        User   `json:"author"`
	PlayURL       string `json:"play_url"`
	CoverURL      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}

const (
	SuccessCode    int32  = 0
	ErrorCode      int32  = 1
	SuccessMessage string = "Success"
)

type FavoriteService interface {
	FavoriteAction(userId int64, videoID int64, actionType int32) (FavoriteActionResponse, error)
	FavoriteList(userID int64) (FavoriteListResponse, error)
	GetUserInfoByIDs(requestingUserID int64, userIDs []int64) ([]*User, error)
	GetFavoriteVideoInfoByUserID(userID int64) ([]*Video, error)
	// GetVideosLikes 批量查询 根据id获取点赞数
	GetVideosLikes(videoIDs []int64) (map[int64]int64, error)
	// AreVideosLikedByUser 查询用户喜欢的视频
	AreVideosLikedByUser(userID int64, videoIDs []int64) (map[int64]bool, error)
}

type FavoriteActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type FavoriteListResponse struct {
	StatusCode int32    `json:"status_code"`
	StatusMsg  string   `json:"status_msg"`
	VideoList  []*Video `json:"video_list"`
}
type FavoriteServiceImpl struct {
}
