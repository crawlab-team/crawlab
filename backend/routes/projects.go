package routes

import (
	"crawlab/constants"
	"crawlab/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

func GetProjectList(c *gin.Context) {
	// 获取列表
	projects, err := model.GetProjectList(nil, 0, "+_id")
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 获取总数
	total, err := model.GetProjectListTotal(nil)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 获取每个项目的爬虫列表
	for i, p := range projects {
		spiders, err := p.GetSpiders()
		if err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
		projects[i].Spiders = spiders
	}

	// 获取未被分配的爬虫数量
	noProject := model.Project{
		Id:          bson.ObjectIdHex(constants.ObjectIdNull),
		Name:        "No Project",
		Description: "Not assigned to any project",
	}
	spiders, err := noProject.GetSpiders()
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	noProject.Spiders = spiders
	projects = append(projects, noProject)

	c.JSON(http.StatusOK, ListResponse{
		Status:  "ok",
		Message: "success",
		Data:    projects,
		Total:   total,
	})
}

func PutProject(c *gin.Context) {
	// 绑定请求数据
	var p model.Project
	if err := c.ShouldBindJSON(&p); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	if err := p.Add(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func PostProject(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
	}

	var item model.Project
	if err := c.ShouldBindJSON(&item); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	if err := model.UpdateProject(bson.ObjectIdHex(id), item); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func DeleteProject(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}

	// 从数据库中删除该爬虫
	if err := model.RemoveProject(bson.ObjectIdHex(id)); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}
