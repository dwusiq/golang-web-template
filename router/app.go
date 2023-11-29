package router

import (
	"github.com/gin-gonic/gin"
	"walletSynV2/handler"
)

// To Java server接口
func initAppRouter(r *gin.RouterGroup) {
	webApiRouter := r.Group("/web")
	{
		//ArcToken api
		webApiRouter.POST("/getBlackList", handler.GetBlackList)   //获取黑名单列表
		webApiRouter.POST("/checkAccount", handler.CheckAccount)   //检查账户是否在黑名单内
		webApiRouter.POST("/deleteAccount", handler.DeleteAccount) //把账户从黑名单中移除

	}
}
