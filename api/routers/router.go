package routers

import (
	"deercoder-chat/api/controllers"
	"deercoder-chat/api/controllers/chat"
	"deercoder-chat/api/controllers/file"
	"github.com/dreamlu/go-tool"
	"github.com/dreamlu/go-tool/tool/result"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"net/http"
	"strconv"
	"strings"
)

func SetRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	router := gin.New()
	gt.MaxUploadMemory = router.MaxMultipartMemory
	router.Use(Cors())

	//router.Use(CheckLogin()) //简单登录验证

	// load the casbin model and policy from files, database is also supported.
	//权限中间件
	//e := casbin.NewEnforcer("conf/authz_model.conf", "conf/authz_policy.csv")
	//router.Use(controllers.NewAuthorizer(e))

	// 接口路由
	// 静态目录
	router.Static("v1/static", "static")

	// Ping test
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	//组的路由,version
	v1 := router.Group("/v1")
	{
		user := controllers.UserService{}
		login := controllers.LoginService{}
		v := v1
		//用户登录
		v.POST("/login/login", login.Login)
		//文件上传
		v.POST("/file/upload", file.UpoadFile)

		// 聊天服务
		chats := v.Group("/chat")
		{
			// websocket get request
			chats.GET("/chatWs", chat.ChatWS)

			// chat message
			chats.POST("/disGroup", chat.DistributeGroup)
			chats.GET("/allMsg", chat.GetAllGroupMsg)
			chats.GET("/lastMsg", chat.GetGroupLastMsg)
			chats.POST("/readLastMsg", chat.ReadGroupLastMsg)
			chats.GET("/getGroupUser", chat.GetGroupUser)
			chats.GET("/getUserList", chat.GetUserList)
			chats.GET("/getUserSearchList", chat.GetUserSearchList)
		}

		// 用户服务
		users := v.Group("/user")
		{
			users.POST("/create", user.Create)
			users.PUT("/update", user.Update)
			users.DELETE("/delete/:id", user.Delete)
			users.GET("/search", user.GetBySearch)
			users.GET("/id/:id", user.GetByID)
		}
	}
	//不存在路由
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": 404,
			"msg":    "接口不存在->('.')/请求方法错误",
		})
	})
	return router
}

// 登录验证
func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 忽略路径
		paths := []string{"/login", "/static"}

		// url
		path := c.Request.URL.String()

		for _, v := range paths {
			if strings.Contains(path, v) {
				return
			}
		}

		// 缓存验证
		var cache gt.CacheManager = new(gt.RedisManager)
		_ = cache.NewCache()
		//cacheModel := der.CacheModel{}
		// redis 存储用户信息
		uid, err := c.Cookie("uid")
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, result.MapNoAuth)
		}
		// key 的类型需要匹配
		// 这里只是简单验证
		// 可做些身份权限等验证
		uidInt64, err := strconv.ParseInt(uid, 10, 64)
		if err != nil {
			log.Info("[路由错误]:", err.Error())
			return
		}
		cacheModel, err := cache.Get(uidInt64)
		if err != nil {
			log.Info("[路由错误]:", err.Error())
			return
		}
		if cacheModel.Data == nil {
			c.Abort()
			c.JSON(http.StatusOK, result.MapNoAuth)
		}
	}
}

// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		//fmt.Println(method)
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法，因为有的模板是要请求两次的
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
