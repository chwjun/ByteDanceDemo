package service

const (
	SuccessCode    int32  = 0
	ErrorCode      int32  = 1
	SuccessMessage string = "Success"
)

type (
	FavoriteService interface {
		FavoriteAction(videoID int64, actionType int32) (FavoriteActionResponse, error)
		FavoriteList(userID int64) (FavoriteListResponse, error)
		GetUserInfoByID(requestingUserID *int64, userID int64) (*User, error)
		GetFavoriteVideoInfoByUserID(userID int64) ([]*Video, error)
	}
)

type FavoriteActionResponse struct {
	StatusCode int32
	StatusMsg  string
}

type FavoriteListResponse struct {
	StatusCode int32
	StatusMsg  string
	VideoList  []*ResponseVideo
}
