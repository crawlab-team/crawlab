package routes

import (
	"crawlab/constants"
	"crawlab/entity"
	"crawlab/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func GetLangList(c *gin.Context) {
	nodeId := c.Param("id")
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    services.GetLangList(nodeId),
	})
}

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

func GetInstalledDepList(c *gin.Context) {
	nodeId := c.Param("id")
	lang := c.Query("lang")
	var depList []entity.Dependency
	if lang == constants.Python {
		if services.IsMasterNode(nodeId) {
			list, err := services.GetPythonLocalInstalledDepList(nodeId)
			if err != nil {
				HandleError(http.StatusInternalServerError, c, err)
				return
			}
			depList = list
		} else {
			list, err := services.GetPythonRemoteInstalledDepList(nodeId)
			if err != nil {
				HandleError(http.StatusInternalServerError, c, err)
				return
			}
			depList = list
		}
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

func GetAllDepList(c *gin.Context) {
	lang := c.Query("lang")
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

func InstallDep(c *gin.Context) {
	type ReqBody struct {
		Lang    string `json:"lang"`
		DepName string `json:"dep_name"`
	}

	nodeId := c.Param("id")

	var reqBody ReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleError(http.StatusBadRequest, c, err)
	}

	if reqBody.Lang == constants.Python {
		if services.IsMasterNode(nodeId) {
			_, err := services.InstallPythonLocalDep(reqBody.DepName)
			if err != nil {
				HandleError(http.StatusInternalServerError, c, err)
				return
			}
		} else {
			_, err := services.InstallPythonRemoteDep(nodeId, reqBody.DepName)
			if err != nil {
				HandleError(http.StatusInternalServerError, c, err)
				return
			}
		}
	} else {
		HandleErrorF(http.StatusBadRequest, c, fmt.Sprintf("%s is not implemented", reqBody.Lang))
		return
	}

	// TODO: check if install is successful

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

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

	if reqBody.Lang == constants.Python {
		if services.IsMasterNode(nodeId) {
			_, err := services.UninstallPythonLocalDep(reqBody.DepName)
			if err != nil {
				HandleError(http.StatusInternalServerError, c, err)
				return
			}
		} else {
			_, err := services.UninstallPythonRemoteDep(nodeId, reqBody.DepName)
			if err != nil {
				HandleError(http.StatusInternalServerError, c, err)
				return
			}
		}
	} else {
		HandleErrorF(http.StatusBadRequest, c, fmt.Sprintf("%s is not implemented", reqBody.Lang))
		return
	}

	// TODO: check if uninstall is successful

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}
