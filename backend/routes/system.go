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
	} else if reqBody.Lang == constants.Nodejs {
		if services.IsMasterNode(nodeId) {
			_, err := services.InstallNodejsLocalDep(reqBody.DepName)
			if err != nil {
				HandleError(http.StatusInternalServerError, c, err)
				return
			}
		} else {
			_, err := services.InstallNodejsRemoteDep(nodeId, reqBody.DepName)
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
	} else if reqBody.Lang == constants.Nodejs {
		if services.IsMasterNode(nodeId) {
			_, err := services.UninstallNodejsLocalDep(reqBody.DepName)
			if err != nil {
				HandleError(http.StatusInternalServerError, c, err)
				return
			}
		} else {
			_, err := services.UninstallNodejsRemoteDep(nodeId, reqBody.DepName)
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
