package routes

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func GetFile(c *gin.Context) {
	path := c.Query("path")
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    string(fileBytes),
	})
}
