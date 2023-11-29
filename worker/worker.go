package worker

import (
	"walletSynV2/etc"
)

type Worker struct {
}

func NewWorker() *Worker {
	return &Worker{}
}

//// 定时任务
//func (w *Worker) timer() {
//	loc, _ := time.LoadLocation("Local")
//	c := cron.NewWithLocation(loc)
//	err := c.AddFunc("0 0/10 * * * ?  ", w.AddContract) // 每天0点过5分运行一次/统计一天的收益
//	if err != nil {
//		zlog.Zlog.Error("start add contract task failed")
//	}
//	c.Start()
//}

func (w *Worker) Run() {
	//go w.timer()
	//go w.UpdateWallet()

	if etc.Conf.Server.ChainMode == 0 {
		// go w.GetEthBlock()
	} else if etc.Conf.Server.ChainMode == 2 {
		go w.GetBtcBlock()
	}
}
