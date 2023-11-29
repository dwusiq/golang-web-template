package router

import (
	"fmt"
	"strings"
	"walletSynV2/middleware"
	"walletSynV2/utils"

	_ "walletSynV2/docs"
	. "walletSynV2/utils/zlog_sing"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init(port int, mode string) error {
	Zlog.Info("starting server...")
	switch mode {
	case gin.DebugMode, gin.ReleaseMode, gin.TestMode:
	default:
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)
	router := gin.Default()
	router.Use(middleware.Cors()) // 跨域处理

	// 文档接口
	doc := router.Group("/docs")
	{
		v1 := doc.Group("/v1")
		{
			v1.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}
	}

	api := router.Group("/api")
	{
		initAppRouter(api)
	}

	address, err := utils.GetIPAddress()
	if err != nil {
		Zlog.Error("Get ip Address fail: " + err.Error())
	}

	if strings.EqualFold(address, "10.0.0.1") {
		address = "127.0.0.1"
	}

	Zlog.Info("server started success...")
	Zlog.Info("swagger url:  " + fmt.Sprintf("http://%s:%d/docs/v1/index.html", address, port))
	return router.Run(fmt.Sprintf("%s:%d", address, port))
}
