package routes

import (
	"crawlab/services"
	"fmt"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"github.com/spf13/viper"
	"net/http"
	"runtime/debug"
)

func GetRepoList(c *gin.Context) {
	var data ListRequestData
	if err := c.ShouldBindQuery(&data); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}
	params := req.Param{
		"page_num":  data.PageNum,
		"page_size": data.PageSize,
		"keyword":   data.Keyword,
		"sort_key":  data.SortKey,
	}
	res, err := req.Get(fmt.Sprintf("%s/public/repos", viper.GetString("repo.apiUrl")), params)
	if err != nil {
		log.Error("get repos error: " + err.Error())
		debug.PrintStack()
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	var resJson interface{}
	if err := res.ToJSON(&resJson); err != nil {
		log.Error("to json error: " + err.Error())
		debug.PrintStack()
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, resJson)
}

func GetRepoSubDirList(c *gin.Context) {
	params := req.Param{
		"full_name": c.Query("full_name"),
	}
	res, err := req.Get(fmt.Sprintf("%s/public/repos/sub-dir", viper.GetString("repo.apiUrl")), params)
	if err != nil {
		log.Error("get repo sub-dir error: " + err.Error())
		debug.PrintStack()
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	var resJson interface{}
	if err := res.ToJSON(&resJson); err != nil {
		log.Error("to json error: " + err.Error())
		debug.PrintStack()
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, resJson)
}

func DownloadRepo(c *gin.Context) {
	type RequestData struct {
		FullName string `json:"full_name"`
	}
	var reqData RequestData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}
	if err := services.DownloadRepo(reqData.FullName, services.GetCurrentUserId(c)); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}
