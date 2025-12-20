package main

/**
 * @author: biao
 * @date: 2025/12/18 下午9:46
 * @description:
 */

import (
	"github.com/bbbbbbbbiao/WeBook/webook/internal/web"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	u := web.NewUserHandler()

	u.RegisterRoutes(server)

	server.Run(":8080")
}
