package scan

import (
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
	"strings"
	"walletSynV2/etc"
	"walletSynV2/service"
	zlog "walletSynV2/utils/zlog_sing"
)

func EthTransactions(block *types.Block) {
	for _, tx := range block.Transactions() {
		if tx.To() == nil { //不是合约调用
			continue
		}

		_, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
		if err != nil {
			zlog.Zlog.Info("get from address fail"+tx.Hash().String(), zap.Error(err))
			continue
		}

		if strings.EqualFold(tx.To().Hex(), etc.Conf.EthContract.ArcTokenAddr.Hex()) &&
			tx.To().Hex() != "0x0000000000000000000000000000000000000000" {
			service.ArcTokenScan(block.Number(), tx)
			continue
		}
	}
}
