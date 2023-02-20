package user_account

import (
	"net/http"

	"github.com/TianYi2000/simple-tiktok/models"
	"github.com/TianYi2000/simple-tiktok/service/user_account"
	"github.com/gin-gonic/gin"
)

type UserLoginResponse struct {
	models.CommonResponse
	*user_account.LoginResponse
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	userLoginResponse, err := user_account.QueryUserLogin(username, password)

	//用户不存在
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			CommonResponse: models.CommonResponse{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	//用户存在
	c.JSON(http.StatusOK, UserLoginResponse{
		CommonResponse: models.CommonResponse{StatusCode: 0},
		LoginResponse:  userLoginResponse,
	})
}
