package routes

import (
	"crawlab/constants"
	"crawlab/model"
	"crawlab/services"
	"crawlab/services/context"
	"crawlab/utils"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type UserListRequestData struct {
	PageNum  int `form:"page_num"`
	PageSize int `form:"page_size"`
}

type UserRequestData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}

// @Summary Get user
// @Description user
// @Tags user
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "user id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	user, err := model.GetUser(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    user,
	})
}

// @Summary Get user list
// @Description Get user list
// @Tags token
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param data body routes.UserListRequestData true "data body"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /users [get]
func GetUserList(c *gin.Context) {
	// 绑定数据
	data := UserListRequestData{}
	if err := c.ShouldBindQuery(&data); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}
	if data.PageNum == 0 {
		data.PageNum = 1
	}
	if data.PageSize == 0 {
		data.PageNum = 10
	}

	// 获取用户列表
	users, err := model.GetUserList(nil, (data.PageNum-1)*data.PageSize, data.PageSize, "-create_ts")
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 获取总用户数
	total, err := model.GetUserListTotal(nil)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 去除密码
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, ListResponse{
		Status:  "ok",
		Message: "success",
		Data:    users,
		Total:   total,
	})
}

// @Summary Put user
// @Description Put user
// @Tags user
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param reqData body routes.UserRequestData true "reqData body"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /users [put]
func PutUser(c *gin.Context) {
	// 绑定请求数据
	var reqData UserRequestData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 默认为正常用户
	if reqData.Role == "" {
		reqData.Role = constants.RoleNormal
	}

	// UserId
	uid := services.GetCurrentUserId(c)

	// 空 UserId 处理
	if uid == "" {
		uid = bson.ObjectIdHex(constants.ObjectIdNull)
	}

	// 添加用户
	if err := services.CreateNewUser(reqData.Username, reqData.Password, reqData.Role, reqData.Email, uid); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Post user
// @Description Post user
// @Tags user
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param item body model.User true "user body"
// @Param id path string true "user id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /users/{id} [post]
func PostUser(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}

	var item model.User
	if err := c.ShouldBindJSON(&item); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	if item.UserId.Hex() == "" {
		item.UserId = bson.ObjectIdHex(constants.ObjectIdNull)
	}

	if err := model.UpdateUser(bson.ObjectIdHex(id), item); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Delete user
// @Description Delete user
// @Tags user
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "user id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}

	// 从数据库中删除该爬虫
	if err := model.RemoveUser(bson.ObjectIdHex(id)); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func Login(c *gin.Context) {
	// 绑定请求数据
	var reqData UserRequestData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		HandleError(http.StatusUnauthorized, c, errors.New("not authorized"))
		return
	}

	// 获取用户
	user, err := model.GetUserByUsername(strings.ToLower(reqData.Username))
	if err != nil {
		HandleError(http.StatusUnauthorized, c, errors.New("not authorized"))
		return
	}

	// 校验密码
	encPassword := utils.EncryptPassword(reqData.Password)
	if user.Password != encPassword {
		HandleError(http.StatusUnauthorized, c, errors.New("not authorized"))
		return
	}

	// 获取token
	tokenStr, err := services.MakeToken(&user)
	if err != nil {
		HandleError(http.StatusUnauthorized, c, errors.New("not authorized"))
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    tokenStr,
	})
}

func GetMe(c *gin.Context) {
	ctx := context.WithGinContext(c)
	user := ctx.User()
	if user == nil {
		ctx.FailedWithError(constants.ErrorUserNotFound, http.StatusUnauthorized)
		return
	}
	ctx.Success(struct {
		*model.User
		Password string `json:"password,omitempty"`
	}{
		User: user,
	}, nil)
}

func PostMe(c *gin.Context) {
	ctx := context.WithGinContext(c)
	user := ctx.User()
	if user == nil {
		ctx.FailedWithError(constants.ErrorUserNotFound, http.StatusUnauthorized)
		return
	}
	var reqBody model.User
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleErrorF(http.StatusBadRequest, c, "invalid request")
		return
	}
	if reqBody.Email != "" {
		user.Email = reqBody.Email
	}
	if reqBody.Setting.NotificationTrigger != "" {
		user.Setting.NotificationTrigger = reqBody.Setting.NotificationTrigger
	}
	if reqBody.Setting.DingTalkRobotWebhook != "" {
		user.Setting.DingTalkRobotWebhook = reqBody.Setting.DingTalkRobotWebhook
	}
	if reqBody.Setting.WechatRobotWebhook != "" {
		user.Setting.WechatRobotWebhook = reqBody.Setting.WechatRobotWebhook
	}
	user.Setting.EnabledNotifications = reqBody.Setting.EnabledNotifications
	user.Setting.ErrorRegexPattern = reqBody.Setting.ErrorRegexPattern
	if reqBody.Setting.MaxErrorLog != 0 {
		user.Setting.MaxErrorLog = reqBody.Setting.MaxErrorLog
	}
	user.Setting.LogExpireDuration = reqBody.Setting.LogExpireDuration

	if user.UserId.Hex() == "" {
		user.UserId = bson.ObjectIdHex(constants.ObjectIdNull)
	}

	if err := user.Save(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func PostMeChangePassword(c *gin.Context) {
	ctx := context.WithGinContext(c)
	user := ctx.User()
	if user == nil {
		ctx.FailedWithError(constants.ErrorUserNotFound, http.StatusUnauthorized)
		return
	}
	var reqBody model.User
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleErrorF(http.StatusBadRequest, c, "invalid request")
		return
	}
	if reqBody.Password == "" {
		HandleErrorF(http.StatusBadRequest, c, "password is empty")
		return
	}
	if user.UserId.Hex() == "" {
		user.UserId = bson.ObjectIdHex(constants.ObjectIdNull)
	}
	user.Password = utils.EncryptPassword(reqBody.Password)
	if err := user.Save(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}
