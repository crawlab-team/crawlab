package routes

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/model"
	"crawlab/services"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

func GetProjectList(c *gin.Context) {
	tag := c.Query("tag")

	// 筛选条件
	query := bson.M{}
	if tag != "" {
		query["tags"] = tag
	}

	// 获取校验
	query = services.GetAuthQuery(query, c)

	// 获取列表
	projects, err := model.GetProjectList(query, "+_id")
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 获取总数
	total, err := model.GetProjectListTotal(query)
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
	if tag == "" {
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
	}

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

	// UserId
	p.UserId = services.GetCurrentUserId(c)

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

	// 获取相关的爬虫
	var spiders []model.Spider
	s, col := database.GetCol("spiders")
	defer s.Close()
	if err := col.Find(bson.M{"project_id": bson.ObjectIdHex(id)}).All(&spiders); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 将爬虫的项目ID置空
	for _, spider := range spiders {
		spider.ProjectId = bson.ObjectIdHex(constants.ObjectIdNull)
		if err := spider.Save(); err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func GetProjectTags(c *gin.Context) {
	type Result struct {
		Tag string `json:"tag" bson:"tag"`
	}

	s, col := database.GetCol("projects")
	defer s.Close()

	pipeline := []bson.M{
		{
			"$unwind": "$tags",
		},
		{
			"$group": bson.M{
				"_id": "$tags",
			},
		},
		{
			"$sort": bson.M{
				"_id": 1,
			},
		},
		{
			"$addFields": bson.M{
				"tag": "$_id",
			},
		},
	}

	var items []Result
	if err := col.Pipe(pipeline).All(&items); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    items,
	})
}
