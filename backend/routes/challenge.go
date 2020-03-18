package routes

import (
	"crawlab/constants"
	"crawlab/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetChallengeList(c *gin.Context) {
	// 获取列表
	users, err := model.GetChallengeList(nil, 0, constants.Infinite, "create_ts")
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 获取总数
	total, err := model.GetChallengeListTotal(nil)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Status:  "ok",
		Message: "success",
		Data:    users,
		Total:   total,
	})
}
