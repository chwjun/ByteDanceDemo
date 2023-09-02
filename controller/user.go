package controller

import (
	"bytedancedemo/model"
	"bytedancedemo/service"
	"bytedancedemo/utils/encryption"
	"bytedancedemo/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User service.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	passwordKey := encryption.Encrypt(password)

	usi := service.GetUserServiceInstance()
	if _, isExist := usi.GetUserBasicByPassword(username, passwordKey); isExist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户已存在"},
		})
	} else {
		user, ok := usi.InsertUser(&model.User{
			Name:            username,
			Password:        passwordKey,
			Role:            "common_user",
			Avatar:          "https://sample-douyin-video.oss-cn-beijing.aliyuncs.com/avatar1.jpeg",
			BackgroundImage: "https://sample-douyin-video.oss-cn-beijing.aliyuncs.com/background.jpeg",
		})
		if !ok {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "注册失败"},
			})
		} else {
			tokenString, err := token.GenerateToken(
				[]byte(viper.GetString("settings.jwt.secretKey")),
				token.Claims{
					UserID:   user.ID,
					UserName: user.Name,
					Role:     user.Role,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour * viper.GetDuration("settings.jwt.expirationTime")).Unix(),
					},
				},
			)
			if err != nil {
				c.JSON(http.StatusOK, UserLoginResponse{
					Response: Response{StatusCode: 1, StatusMsg: "token令牌签发失败"},
				})
			} else {
				c.JSON(http.StatusOK, UserLoginResponse{
					Response: Response{StatusCode: 0},
					UserId:   user.ID,
					Token:    tokenString,
				})
			}
		}
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	passwordKey := encryption.Encrypt(password)

	usi := service.GetUserServiceInstance()
	user, isExist := usi.GetUserBasicByPassword(username, passwordKey)
	if !isExist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "用户不存在",
			},
		})
	} else {
		tokenString, err := token.GenerateToken([]byte(viper.GetString("settings.jwt.secretKey")), token.Claims{
			UserID:   user.ID,
			UserName: user.Name,
			Role:     user.Role,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * viper.GetDuration("settings.jwt.expirationTime")).Unix(),
			},
		})
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "token令牌签发失败",
				},
			})
		} else {
			zap.L().Debug("登录成功", zap.Int64("userID", user.ID), zap.String("username", user.Name), zap.String("role", user.Role))
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{
					StatusCode: 0,
				},
				UserId: user.ID,
				Token:  tokenString,
			})
		}
	}
}

func UserInfo(c *gin.Context) {
	userID := c.GetInt64("user_id")
	usi := service.GetUserServiceInstance()
	user, err := usi.GetUserDetailsById(userID, nil)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
	} else {
		zap.L().Debug("查询用户详情成功", zap.Any("user", user))
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     *user,
		})
	}
}
