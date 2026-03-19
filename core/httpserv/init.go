package httpserv

import (
	"fmt"
	"gServ/core/config"
	"gServ/core/log"
	"gServ/pkg/gserv"
	"gServ/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine

	captcha_generator          = gserv.NewCaptchaGenerator(gserv.CAPTCHA_LENGTH)
	captcha_template_generator = gserv.NewCaptchaTemplateGenerator()
)

func Init() error {
	// 设置运行模式
	switch config.GetConfig().Server.Mode {
	case config.SERVER_MODE_DEV:
		gin.SetMode(gin.DebugMode)
	case config.SERVER_MODE_PROD:
		gin.SetMode(gin.ReleaseMode)
	default:
		return fmt.Errorf("无效服务模式: %s", config.GetConfig().Server.Mode)
	}

	// 创建路由
	router = gin.New()

	initRouter(router)

	err := captcha_template_generator.Open(config.GetConfig().Server.Email.Template)
	if err != nil {
		return err
	}

	return nil
}

func initRouter(router *gin.Engine) {
	// 初始化通用中间件
	router.
		Use(middleware.ZapLogger(log.GetZapLogger())).
		Use(middleware.StderrLogger(log.GetStdErrorLogger())).
		Use(middleware.StdoutLogger(log.GetStdInfoLogger())).
		Use(gin.Recovery()).
		Use(middleware.Cors())

	// 接口路由
	api_router := router.Group("/api")
	{
		// 健康检查
		api_router.GET("/health", get_Api_Health)

		games_router := api_router.Group("/games")
		games_router.Use(middleware.PlayerAuth())
		{
			// 获取游戏列表
			games_router.GET("", get_Api_Games)
		}

		rooms_router := api_router.Group("/rooms")
		rooms_router.Use(middleware.PlayerAuth())
		{
			// 获取房间列表，返回房间ID、房间名、房间最大人数、房间当前人数的列表和房间创建时间
			rooms_router.GET("", get_Api_Rooms)
		}

		// 房间相关
		room_router := api_router.Group("/room")
		room_router.Use(middleware.PlayerAuth())
		{
			// 创建房间，创建后不会主动加入房间，还需要再调用一次TCP服务建立连接后加入房间
			room_router.POST("", post_Api_Room)
			// 获取房间信息，返回房间ID、房间名、房间最大人数、房间当前人数、房间其他玩家信息列表和房间创建时间
			room_router.GET("/:game_id/:room_id", get_Api_Room)
			// 放逐玩家
			room_router.PUT("/:game_id/:room_id/:player_id", put_Api_Room_ExilePlayer)
			// 锁定房间
			room_router.PUT("/:game_id/:room_id/lock", put_Api_Room_Lock)
			// 解锁房间
			room_router.PUT("/:game_id/:room_id/unlock", put_Api_Room_Unlock)
			// 删除房间，需要玩家权限为房主
			room_router.DELETE("/:game_id/:room_id", delete_Api_Room)
		}

		// 验证码相关
		captcha_router := api_router.Group("/captcha")
		{
			// 服务端向指定邮箱发送一个验证码，请求需要有邮箱账号、验证码类型
			captcha_router.POST("/email", post_Api_Captcha_Email)
		}

		// 玩家相关
		player_router := api_router.Group("/player")
		{
			// 玩家注册，需要注册方式（邮箱注册）、邮箱、密码和验证码
			player_router.POST("/register", post_Api_Player_Register)
			// 玩家登录，输入玩家ID和密码，返回token
			player_router.POST("/login", post_Api_Player_Login)

			// 玩家权限验证，所有和ID相关的操作首先要把token中的ID和URL中的ID进行对比，如果一致则继续，不一致则返回错误
			player_router.Use(middleware.PlayerAuth())
			{
				// 获取玩家信息，返回玩家ID、昵称、邮箱、当前所在房间和账号创建时间
				player_router.GET("/:player_id", get_Api_Player)
				// 更新玩家信息，目前仅需支持修改昵称和邮箱，修改邮箱需要附带发送到旧邮箱的验证码
				player_router.PUT("", put_Api_Player)
				// 修改密码，需要输入旧密码和新密码还有发送到邮箱的验证码
				player_router.PUT("/password", put_Api_Player_Password)
				// 删除玩家
				player_router.DELETE("/:player_id", delete_Api_Player)
			}
		}

		// 数据存储相关 - JSON数据增删查改接口
		data_router := api_router.Group("/data")
		data_router.Use(middleware.PlayerAuth())
		{
			// 判断存档是否存在，不存在开始游戏前，可以调用创建存档接口进行初始化
			data_router.GET("/:game_id/exists", get_Api_Data_Exists)
			// 创建JSON数据
			data_router.POST("", post_Api_Data)
			// 获取JSON数据
			data_router.GET("", get_Api_Data)
			// 更新JSON数据
			data_router.PUT("/:game_id", put_Api_Data)
			// 删除JSON数据
			data_router.DELETE("", delete_Api_Data)
		}
	}

	// 管理路由，所有请求头里面只需要有AuthCode字段进行验证即可
	admin_router := router.Group("/admin")
	admin_router.Use(middleware.CodeAuth())
	{
		games_router := admin_router.Group("/games")
		{
			// 获取游戏列表
			games_router.GET("", get_Admin_Games)
		}

		game_router := admin_router.Group("/game")
		{
			// 创建游戏
			game_router.POST("", post_Admin_Game)
			// 删除游戏
			game_router.DELETE("/:game_id", delete_Admin_Game)
		}

		room_router := admin_router.Group("/room")
		{
			// 强制删除房间
			room_router.DELETE("/:game_id/:room_id", delete_Admin_Room)
		}

		ban_router := admin_router.Group("/ban")
		{
			player_router := ban_router.Group("/player")
			{
				// 封禁玩家，直接在数据库中Delete玩家即可
				player_router.DELETE("/:player_id", delete_Admin_Ban_Player)
				// 解封玩家，将玩家在数据库中的DeletedAt字段置为空即可
				player_router.PUT("/:player_id", put_Admin_Ban_Player)
			}
		}
	}
}

func GetHTTPServer() *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", config.GetConfig().Server.HTTPPort),
		Handler: router,
	}
}
