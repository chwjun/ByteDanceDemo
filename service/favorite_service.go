package service

type User struct {
	ID              int64
	Name            string
	FollowCount     int64
	FollowerCount   int64
	IsFollow        bool
	Avatar          string
	BackgroundImage string
	Signature       string
	TotalFavorited  int64
	WorkCount       int64
	FavoriteCount   int64
}

type Video struct {
	ID            int64
	Author        User
	PlayURL       string
	CoverURL      string
	FavoriteCount int64
	CommentCount  int64
	IsFavorite    bool
	Title         string
}

const (
	SuccessCode    int32  = 0
	ErrorCode      int32  = 1
	SuccessMessage string = "Success"
)

type FavoriteService interface {
	FavoriteAction(videoID int64, actionType int32) (FavoriteActionResponse, error)
	FavoriteList(userID int64) (FavoriteListResponse, error)
	GetUserInfoByID(requestingUserID *int64, userID int64) (*User, error)
	GetFavoriteVideoInfoByUserID(userID int64) ([]*Video, error)
}

type FavoriteActionResponse struct {
	StatusCode int32
	StatusMsg  string
}

type FavoriteListResponse struct {
	StatusCode int32
	StatusMsg  string
	VideoList  []*Video
}
type FavoriteServiceImpl struct {
}
