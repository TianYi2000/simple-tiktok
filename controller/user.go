package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	MaxUsernameLength = 32
	MaxPasswordLength = 32
	MinPasswordLength = 6
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

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

// Register funtions

func checkAccountLen(username, password string) error {
	// println(len(username), len(password))
	if username == "" {
		return errors.New("用户名为空")
	}
	if len(username) > MaxUsernameLength {
		return errors.New("用户名长度超出限制")
	}
	if password == "" {
		return errors.New("密码为空")
	}
	if len(password) > MaxPasswordLength {
		return errors.New("密码长度超出限制")
	}
	if len(password) < MinPasswordLength {
		return errors.New("密码长度过短")
	}
	return nil
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")

	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
