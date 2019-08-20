package routes

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
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
