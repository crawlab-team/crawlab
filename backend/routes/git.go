package routes

import (
	"crawlab/model"
	"crawlab/services"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

func GetGitRemoteBranches(c *gin.Context) {
	url := c.Query("url")
	branches, err := services.GetGitRemoteBranches(url)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    branches,
	})
}

func GetGitSshPublicKey(c *gin.Context) {
	content := services.GetGitSshPublicKey()
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    content,
	})
}

func GetGitCommits(c *gin.Context) {
	spiderId := c.Query("spider_id")
	if spiderId == "" || !bson.IsObjectIdHex(spiderId) {
		HandleErrorF(http.StatusInternalServerError, c, "invalid request")
		return
	}
	spider, err := model.GetSpider(bson.ObjectIdHex(spiderId))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	commits, err := services.GetGitCommits(spider)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    commits,
	})
}
