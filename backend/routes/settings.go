package routes

import (
	"crawlab/services/context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func GetSettings(c *gin.Context) {
	ctx := context.WithGinContext(c)

	currentUser := ctx.User()

	//enable register
	canRegister := viper.GetBool("server.master_settings,auth.register.enable")

	//auth
	if currentUser != nil {

	}

	ctx.Success(gin.H{
		"can_register": canRegister,
	})
}
