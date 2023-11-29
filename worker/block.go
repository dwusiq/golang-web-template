package worker

import (
	"context"
	"math/big"
	"time"
	"walletSynV2/etc"
	"walletSynV2/models"
	"walletSynV2/node"
	zlog "walletSynV2/utils/zlog_sing"
	"walletSynV2/worker/scan"

	"go.uber.org/zap"
)

func (w *Worker) GetEthBlock() {
	currentHeight, err := models.ReadChainHeight()
	if err != nil {
		zlog.Zlog.Error("read chain height err", zap.Error(err))
	}

	for {
		oldHight := currentHeight

		if currentHeight < etc.Conf.Server.StartHeight {
			currentHeight = etc.Conf.Server.StartHeight
		}
		zlog.Zlog.Info("synchronizing block currentHeight: ", zap.Int64("currentHeight", currentHeight))
		time.Sleep(1 * time.Millisecond)

		blockNumber := big.NewInt(currentHeight)
		block, err := node.EthClient.GetClient().BlockByNumber(context.Background(), blockNumber)

		if err != nil {
			zlog.Zlog.Info("The maximum block has been reached...")
			zlog.Zlog.Error("err :", zap.Error(err))
			time.Sleep(3000 * time.Millisecond)
			continue
		}

		// scan transactions
		scan.EthTransactions(block)

		currentHeight++
		if err := models.SaveChainHeight(oldHight, currentHeight); err != nil {
			zlog.Zlog.Warn("update chain height err height", zap.Int64("currentHeight", currentHeight),
				zap.Error(err))
		}
	}
}
