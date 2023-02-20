package user_account

import (
	"net/http"

	"github.com/TianYi2000/simple-tiktok/models"
	"github.com/gin-gonic/gin"
)

type UserRegisterResponse struct {
	models.CommonResponse
	*AccountResponse
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	registerResponse, err := AccountHandler(username, password)

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
		CommonResponse:  models.CommonResponse{StatusCode: 0},
		AccountResponse: registerResponse,
	})
}
