package main

import (
	"os"
	"strings"
	"walletSynV2/database"
	"walletSynV2/etc"
	"walletSynV2/models"
	"walletSynV2/node"
	"walletSynV2/router"
	. "walletSynV2/utils/zlog_sing"
	"walletSynV2/worker"

	"go.uber.org/zap"
)

//	@title			合约端程序
//	@version		1.0
//	@description	合约端API文档

//	@contact.name	API Support
//	@contact.author
//	@contact.email

// @host		127.0.0.1:8080
// @BasePath	/api
func main() {
	err := etc.InitConfig("./conf/config.yaml")
	if err != nil {
		Zlog.Error("init conf failed", zap.Error(err))
		os.Exit(1)
	}

	// 初始化节点
	if etc.Conf.Server.ChainMode == 0 {
		err := node.InitEthClient()
		if err != nil {
			Zlog.Error("init eth client failed", zap.Error(err))
			os.Exit(1)
		}
		etc.Conf.Level.ChainId = node.EthClient.GetChainId()

	} else if etc.Conf.Server.ChainMode == 2 {
		err := node.InitBtcClient()
		if err != nil {
			Zlog.Error("init btc client failed", zap.Error(err))
			os.Exit(1)
		}
		// etc.Conf.Level.ChainId = node.EthClient.GetChainId()

	}

	// 初始化Mysql对象模型
	models.RegisterMySQLModel()

	// 初始化Mysql
	isDebug := strings.EqualFold(etc.Conf.Server.RunMode, "debug")
	err = database.InitMysql(etc.Conf.Mysql, isDebug)
	if err != nil {
		Zlog.Error("init MySQL failed", zap.Error(err))
		os.Exit(1)
	}

	// // 初始化LevelDB
	// err = database.InitLevelDB(etc.Conf.Level.ChainId)
	// if err != nil {
	// 	Zlog.Error("init LevelDB failed", zap.Error(err))
	// 	os.Exit(1)
	// }

	// 开始扫块
	worker.NewWorker().Run()

	// 初始化gin router
	err = router.Init(etc.Conf.Server.HttpPort, etc.Conf.Server.RunMode)
	if err != nil {
		Zlog.Error("start server failed", zap.Error(err))
		os.Exit(1)
	}

}
