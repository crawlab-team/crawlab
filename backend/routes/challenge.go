package routes

import (
	"crawlab/constants"
	"crawlab/model"
	"crawlab/services"
	"crawlab/services/challenge"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetChallengeList(c *gin.Context) {
	// 获取列表
	users, err := model.GetChallengeListWithAchieved(nil, 0, constants.Infinite, "create_ts", services.GetCurrentUserId(c))
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

func CheckChallengeList(c *gin.Context) {
	uid := services.GetCurrentUserId(c)
	if err := challenge.CheckChallengeAndUpdateAll(uid); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}
