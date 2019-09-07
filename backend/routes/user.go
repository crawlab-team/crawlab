package routes

import (
	"crawlab/constants"
	"crawlab/model"
	"crawlab/services"
	"crawlab/services/context"
	"crawlab/utils"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"net/http"
	"strings"
	"time"
)

type UserListRequestData struct {
	PageNum  int `form:"page_num"`
	PageSize int `form:"page_size"`
}

type UserRequestData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetUser(c *gin.Context) {
	id := c.Param("id")

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

func PutUser(c *gin.Context) {
	ctx := context.WithGinContext(c)
	// 绑定请求数据
	var reqData struct {
		Username        string `json:"username" binding:"required,min=5"`
		Password        string `json:"password" binding:"required,min=5"`
		ConfirmPassword string `json:"confirm_password" binding:"eqfield=Password"`
	}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		ctx.Failed(err)
		return
	}
	salt := utils.RandomString(10)
	// 添加用户
	user := &model.User{
		Username: strings.ToLower(reqData.Username),
		Role:     constants.RoleNormal,
		Enable:   true,
		Salt:     salt,
		Password: utils.EncryptPasswordV2(reqData.Password, salt),
	}
	if err := user.Add(); err != nil {
		if mgo.IsDup(err) {
			ctx.Failed(constants.ErrorAccountHasExisted)
			return
		}
		ctx.Failed(err)
		return
	}
	ctx.Success(nil)
}

func PostUser(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
	}

	var item model.User
	if err := c.ShouldBindJSON(&item); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
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
	var reqData struct {
		Username string `json:"username" binding:"required,min=5"`
		Password string `json:"password" binding:"required,min=5"`
	}
	ctx := context.WithGinContext(c)
	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.Failed(err)
		return
	}

	// 获取用户
	user, err := model.GetUserByUsername(strings.ToLower(reqData.Username))
	if err != nil {
		if err == mgo.ErrNotFound {
			ctx.Failed(constants.ErrorUsernameOrPasswordInvalid)
			return
		} else {
			ctx.Failed(constants.ErrorMongoError)
			return
		}
	}
	// 校验密码
	isLoggedIn := user.LoginWithPassword(reqData.Password)
	if !isLoggedIn {
		ctx.Failed(constants.ErrorUsernameOrPasswordInvalid)
		return
	}
	if !user.Enable {
		ctx.Failed(constants.ErrorAccountDisabled)
		return
	}
	//if user.RePasswordTs.Before(time.Now()){
	//	ctx.Failed(constants.ErrorNeedResetPassword)
	//	return
	//}
	// 获取token
	tokenStr, err := services.MakeToken(&user)
	if err != nil {
		ctx.Failed(constants.ErrorUsernameOrPasswordInvalid)
		return
	}

	ctx.Success(gin.H{
		"token":          tokenStr,
		"reset_password": user.RePasswordTs.Before(time.Now()),
	})
}
func GetMe(c *gin.Context) {
	ctx := context.WithGinContext(c)
	user := ctx.User()
	if user == nil {
		ctx.Failed(constants.ErrorUserNotFound)
		return
	}
	ctx.Success(struct {
		*model.User
		Password string `json:"password,omitempty"`
	}{
		User: user,
	}, nil)
}
func ChangePassword(c *gin.Context) {
	ctx := context.WithGinContext(c)
	var requestData struct {
		OldPassword        string `json:"old_password" binding:"required,min=5"`
		NewPassword        string `json:"new_password" binding:"required,min=5,nefield=OldPassword"`
		NewConfirmPassword string `json:"confirm_new_password" binding:"required,min=5,eqfield=NewPassword"`
	}
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.Failed(err)
		return
	}
	currentUser := ctx.User()
	if currentUser == nil {
		ctx.Failed(constants.ErrorUserNotFound)
		return
	}
	if !currentUser.ValidatePassword(requestData.OldPassword) {
		ctx.Failed(constants.ErrorUsernameOrPasswordInvalid)
		return
	}

	if err := model.ChangePassword(currentUser, requestData.NewPassword); err != nil {
		ctx.Failed(constants.ErrorMongoError)
		return
	}
	ctx.Success(nil)
}
