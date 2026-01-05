package main

/**
 * @author: biao
 * @date: 2025/12/18 下午9:46
 * @description:
 */

import (
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/repository/dao"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/service"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/web"
	"github.com/bbbbbbbbiao/WeBook/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

func main() {
	// 初始化数据库
	db := initDB()

	server := initWbeServer()

	u := initUser(db)
	u.RegisterRoutes(server)

	server.Run(":8080")
}

func initWbeServer() *gin.Engine {
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

	// 中间件-提取session
	//store := cookie.NewStore([]byte("secret"))
	// 第一个key 身份认证，第二个key 密码加密 (不适应于多实例部署)
	//store := memstore.NewStore([]byte("3E7QYaUxM5tMhDWwd5HphdYWND7WR2Vx"), []byte("Aj6R5sfYMCsxMwsb5TSUjP3228PBdXCE"))
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "", "", []byte("3E7QYaUxM5tMhDWwd5HphdYWND7WR2Vx"), []byte("Aj6R5sfYMCsxMwsb5TSUjP3228PBdXCE"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("my_session", store))

	// 中间件-校验登录
	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePath("/users/login").
		IgnorePath("/users/signup").
		Build())

	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDao(db)
	// 初始化用户模块的Repository
	ur := repository.NewUserRepository(ud)
	// 初始化用户模块的Service
	svc := service.NewUserService(ur)
	// 初始化用户模块的Handler
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		// 表示该goroutine直接退出
		panic(err)
	}
	// 初始化数据库表
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
