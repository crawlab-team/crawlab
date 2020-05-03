package routes

import (
	"crawlab/constants"
	"crawlab/entity"
	"crawlab/services"
	"crawlab/services/rpc"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// @Summary Get language list
// @Description Get language list
// @Tags system
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "node id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /nodes/{id}/langs [get]
func GetLangList(c *gin.Context) {
	nodeId := c.Param("id")
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    services.GetLangList(nodeId),
	})
}

// @Summary Get dep list
// @Description Get dep list
// @Tags system
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "node id"
// @Param lang query string true "language"
// @Param dep_name query string true "dep name"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /nodes/{id}/deps [get]
func GetDepList(c *gin.Context) {
	nodeId := c.Param("id")
	lang := c.Query("lang")
	depName := c.Query("dep_name")

	var depList []entity.Dependency
	if lang == constants.Python {
		list, err := services.GetPythonDepList(nodeId, depName)
		if err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
		depList = list
	} else if lang == constants.Nodejs {
		list, err := services.GetNodejsDepList(nodeId, depName)
		if err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
		depList = list
	} else {
		HandleErrorF(http.StatusBadRequest, c, fmt.Sprintf("%s is not implemented", lang))
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    depList,
	})
}

// @Summary Get installed dep list
// @Description Get installed dep list
// @Tags system
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "node id"
// @Param lang query string true "language"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /nodes/{id}/deps/installed [get]
func GetInstalledDepList(c *gin.Context) {
	nodeId := c.Param("id")
	lang := c.Query("lang")
	var depList []entity.Dependency
	if services.IsMasterNode(nodeId) {
		list, err := rpc.GetInstalledDepsLocal(lang)
		if err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
		depList = list
	} else {
		list, err := rpc.GetInstalledDepsRemote(nodeId, lang)
		if err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
		depList = list
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    depList,
	})
}

// @Summary Get all dep list
// @Description Get all dep list
// @Tags system
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param lang path string true "language"
// @Param dep_nane query string true "dep name"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /system/deps/:lang [get]
func GetAllDepList(c *gin.Context) {
	lang := c.Param("lang")
	depName := c.Query("dep_name")

	// 获取所有依赖列表
	var list []string
	if lang == constants.Python {
		_list, err := services.GetPythonDepListFromRedis()
		if err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
		list = _list
	} else {
		HandleErrorF(http.StatusBadRequest, c, fmt.Sprintf("%s is not implemented", lang))
		return
	}

	// 过滤依赖列表
	var depList []string
	for _, name := range list {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(depName)) {
			depList = append(depList, name)
		}
	}

	// 只取前20
	var returnList []string
	for i, name := range depList {
		if i >= 10 {
			break
		}
		returnList = append(returnList, name)
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    returnList,
	})
}

// @Summary Install  dep
// @Description Install dep
// @Tags system
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "node id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /nodes/{id}/deps/install [Post]
func InstallDep(c *gin.Context) {
	type ReqBody struct {
		Lang    string `json:"lang"`
		DepName string `json:"dep_name"`
	}

	nodeId := c.Param("id")

	var reqBody ReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	if services.IsMasterNode(nodeId) {
		if err := rpc.InstallDepLocal(reqBody.Lang, reqBody.DepName); err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	} else {
		if err := rpc.InstallDepRemote(nodeId, reqBody.Lang, reqBody.DepName); err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Uninstall  dep
// @Description Uninstall dep
// @Tags system
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "node id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /nodes/{id}/deps/uninstall [Post]
func UninstallDep(c *gin.Context) {
	type ReqBody struct {
		Lang    string `json:"lang"`
		DepName string `json:"dep_name"`
	}

	nodeId := c.Param("id")

	var reqBody ReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleError(http.StatusBadRequest, c, err)
	}

	if services.IsMasterNode(nodeId) {
		if err := rpc.UninstallDepLocal(reqBody.Lang, reqBody.DepName); err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	} else {
		if err := rpc.UninstallDepRemote(nodeId, reqBody.Lang, reqBody.DepName); err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Get dep json
// @Description Get dep json
// @Tags system
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param lang path string true "language"
// @Param dep_name path string true "dep name"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /system/deps/{lang}/{dep_name}/json [get]
func GetDepJson(c *gin.Context) {
	depName := c.Param("dep_name")
	lang := c.Param("lang")

	var dep entity.Dependency
	if lang == constants.Python {
		_dep, err := services.FetchPythonDepInfo(depName)
		if err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
		dep = _dep
	} else {
		HandleErrorF(http.StatusBadRequest, c, fmt.Sprintf("%s is not implemented", lang))
		return
	}

	c.Header("Cache-Control", "max-age=86400")
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    dep,
	})
}

// @Summary Install language
// @Description Install language
// @Tags system
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "node id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /nodes/{id}/langs/install [Post]
func InstallLang(c *gin.Context) {
	type ReqBody struct {
		Lang string `json:"lang"`
	}

	nodeId := c.Param("id")

	var reqBody ReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	if services.IsMasterNode(nodeId) {
		_, err := rpc.InstallLangLocal(reqBody.Lang)
		if err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	} else {
		_, err := rpc.InstallLangRemote(nodeId, reqBody.Lang)
		if err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	}

	// TODO: check if install is successful

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}
