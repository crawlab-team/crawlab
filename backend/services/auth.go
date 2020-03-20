package services

import (
	"crawlab/constants"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func GetAuthQuery(query bson.M, c *gin.Context) bson.M {
	user := GetCurrentUser(c)
	if user.Role == constants.RoleAdmin {
		// 获得所有数据
		return query
	} else {
		// 只获取自己的数据
		query["user_id"] = user.Id
		return query
	}
}

