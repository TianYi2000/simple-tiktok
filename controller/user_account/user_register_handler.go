package user_account

import (
	"net/http"

	"github.com/TianYi2000/simple-tiktok/models"
	user_login "github.com/TianYi2000/simple-tiktok/service/user_account"
	"github.com/gin-gonic/gin"
)

type UserRegisterResponse struct {
	models.CommonResponse
	*user_login.LoginResponse
}

func UserRegisterHandler(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// rawVal, _ := c.Get("password")
	// password, ok := rawVal.(string)

	// if !ok {
	// 	c.JSON(http.StatusOK, UserRegisterResponse{
	// 		CommonResponse: models.CommonResponse{
	// 			StatusCode: 1,
	// 			StatusMsg:  "密码解析出错",
	// 		},
	// 	})
	// 	return
	// }
	registerResponse, err := user_login.PostUserLogin(username, password)

	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			CommonResponse: models.CommonResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, UserRegisterResponse{
		CommonResponse: models.CommonResponse{StatusCode: 0},
		LoginResponse:  registerResponse,
	})
}
