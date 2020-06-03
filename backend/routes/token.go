package routes

import (
	"crawlab/model"
	"crawlab/services"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
	"time"
)

// @Summary Get token
// @Description token
// @Tags token
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /tokens [get]
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

// @Summary Put token
// @Description token
// @Tags token
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /tokens [put]
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

// @Summary Delete token
// @Description Delete token
// @Tags token
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "token id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /tokens/{id} [delete]
func DeleteToken(c *gin.Context) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	if err := model.DeleteTokenById(bson.ObjectIdHex(id)); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}
