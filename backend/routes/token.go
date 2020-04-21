package routes

import (
	"crawlab/model"
	"crawlab/services"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
	"time"
)

func GetTokens(c *gin.Context) {
	u := services.GetCurrentUser(c)

	tokens, err := model.GetTokensByUserId(u.Id)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    tokens,
	})
}

func PutToken(c *gin.Context) {
	u := services.GetCurrentUser(c)

	tokenStr, err := services.MakeToken(u)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	t := model.Token{
		Id:       bson.NewObjectId(),
		Token:    tokenStr,
		UserId:   u.Id,
		CreateTs: time.Now(),
		UpdateTs: time.Now(),
	}

	if err := t.Add(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func DeleteToken(c *gin.Context) {
	id := c.Param("id")

	if err := model.DeleteTokenById(bson.ObjectIdHex(id)); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}
