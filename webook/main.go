package main

/**
 * @author: biao
 * @date: 2025/12/18 下午9:46
 * @description:
 */

import (
	"github.com/bbbbbbbbiao/WeBook/webook/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
)

func main() {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:3000"},
		//AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 不填则是全部方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.HasPrefix(origin, "youCompany.com")
		},
	}))
	u := web.NewUserHandler()

	u.RegisterRoutes(server)

	server.Run(":8080")
}
