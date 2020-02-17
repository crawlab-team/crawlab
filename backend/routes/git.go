package routes

import (
	"crawlab/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetGitBranches(c *gin.Context) {
	url := c.Query("url")
	branches, err := services.GetGitBranches(url)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    branches,
	})
}

