package user_account

import (
	"net/http"

	"github.com/TianYi2000/simple-tiktok/models"
	"github.com/gin-gonic/gin"
)

type UserLoginResponse struct {
	models.CommonResponse
	*AccountResponse
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	userLoginResponse, err := QueryUserLogin(username, password)

	//用户不存在
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			CommonResponse: models.CommonResponse{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	//用户存在
	c.JSON(http.StatusOK, UserLoginResponse{
		CommonResponse:  models.CommonResponse{StatusCode: 0},
		AccountResponse: userLoginResponse,
	})
}
