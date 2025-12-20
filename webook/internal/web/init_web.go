package web

import "github.com/gin-gonic/gin"

/**
 * @author: biao
 * @date: 2025/12/19 下午2:49
 * @description: 初始化web层，注册路由
 */

func RegisterRouter() *gin.Engine {
	server := gin.Default()

	u := &UserHandler{}
	server.POST("/users/signup", u.SingUp)
	server.POST("/users/login", u.Login)
	server.POST("/users/edit", u.Edit)
	server.GET("/users/profile", u.Profile)

	return server
}
