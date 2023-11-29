package worker

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"walletSynV2/constant"
	"walletSynV2/etc"
	"walletSynV2/models"
	"walletSynV2/node"
	zlog "walletSynV2/utils/zlog_sing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"go.uber.org/zap"
)

type Inscription struct {
	Protocol  string `json:"p"`
	Operation string `json:"op"`
	TickName  string `json:"tick"`
	MintMax   string `json:"max"`
	MintLimit string `json:"lim"`
	Amount    string `json:"amt"`
}

func (w *Worker) GetBtcBlock() {
	// testMethod()
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

		// 获取指定高度的区块哈希
		blockHash, err := node.BtcClient.GetClient().GetBlockHash(currentHeight)
		if err != nil {
			zap.Error(err)
		}

		// 获取区块信息
		block, err := node.BtcClient.GetClient().GetBlockVerbose(blockHash)
		if err != nil {
			log.Fatal(err)
		}

		// 遍历区块中的每个交易
		for _, txID := range block.Tx {
			fmt.Printf("txID:%s\n", txID)

			// 将交易哈希转换为*chainhash.Hash类型
			txID := "b61b0172d95e266c18aea0c624db987e971a5d6d4ebc2aaed85da4642d635735" //deploy
			// txID := "3a86a15f22572019867f3d551289cce83abaccd28ffffc881364379b8ecdc908"  //transfer
			hash, err := chainhash.NewHashFromStr(txID)
			if err != nil {
				log.Fatal(err)
			}

			// 获取交易详细信息
			tx, err := node.BtcClient.GetClient().GetRawTransactionVerbose(hash)
			if err != nil {
				log.Fatal(err)
			}

			//不解析coinbase交易
			if len(tx.Vin) == 0 || tx.Vin[0].IsCoinBase() {
				zlog.Zlog.Info("coinbase transaction:", zap.String("TxHash", txID))
				continue
			}

			//铭文记录
			inscriptionRsp, err := getInscription(tx)
			if err != nil {
				log.Fatal(err)
			}
			//铭文转JSON   {"p":"brc-20","op":"transfer","tick":"ordi","amt":"10"}
			var inscription Inscription
			err = json.Unmarshal([]byte(inscriptionRsp), &inscription)
			if err != nil {
				fmt.Println("解析JSON失败：", err)
				return
			}
			//必须是目标的协议类型才记录
			protocolType := constant.ProtocolTypeBRC20
			if !strings.EqualFold(inscription.Protocol, "brc-20") {
				continue
			}
			//匹配操作类型
			opType := constant.OpTypeDeploy
			if strings.EqualFold(inscription.Operation, "deploy") {
				opType = constant.OpTypeMint
			} else if strings.EqualFold(inscription.Operation, "mint") {
				opType = constant.OpTypeDeploy
			} else if strings.EqualFold(inscription.Operation, "transfer") {
				opType = constant.OpTypeTransfer
			} else {
				log.Fatal("invalid operate type")
			}

			//查询sender
			sender, err := getSenderAddress(tx)
			if err != nil {
				log.Fatal(err)
			}

			//获取和保存sender
			senderInfo, err := models.FindOrSaveAddress(sender)
			if err != nil {
				log.Fatal(err)
			}

			//获取和记录所有输出用户
			receiptList, _, err := getAllReceiptAddressAndAmount(tx)
			if err != nil {
				log.Fatal(err)
			}

			mintMaxInt, mintLimitInt := int64(0), int64(0)
			if inscription.MintMax != "" {
				mintMaxInt, err = strconv.ParseInt(inscription.MintMax, 10, 64)
				if err != nil {
					fmt.Println("无法将字符串转换为整数:", err)
					return
				}
			}

			if inscription.MintLimit != "" {
				mintLimitInt, err = strconv.ParseInt(inscription.MintLimit, 10, 64)
				if err != nil {
					fmt.Println("无法将字符串转换为整数:", err)
					return
				}
			}

			protocolInfo := &models.ProtocolInfo{
				Deployer:      senderInfo.Id,
				DeployedTxId:  txID,
				DeployedBlock: currentHeight,
				ProtocolType:  protocolType,
				TickName:      inscription.TickName,
				ProtocolJson:  inscriptionRsp,
				MintMax:       mintMaxInt,
				MintLimit:     mintLimitInt,
			}
			protocolId, err := models.SaveProtocol(protocolInfo)
			if err != nil {
				log.Fatal(err)
			}

			transferToId := senderInfo.Id
			for i, receipt := range receiptList {
				receiptInfo, err := models.FindOrSaveAddress(receipt)
				if err != nil {
					log.Fatal(err)
				}
				if !strings.EqualFold(receipt, sender) && i == 0 {
					transferToId = receiptInfo.Id
				}
			}

			//保存操作记录
			ticketAmountFloat := float64(0)
			if inscription.Amount != "" {
				ticketAmountFloat, err = strconv.ParseFloat(inscription.Amount, 64)
				if err != nil {
					log.Fatal(err)
				}
			}

			opHistory := models.OperationHistory{
				ProtocolId:    protocolId,
				OperationType: opType,
				OperatedTxId:  txID,
				OperatedBlock: currentHeight,
				FromUser:      senderInfo.Id,
				ToUser:        transferToId,
				Amount:        ticketAmountFloat,
			}

			err = models.SaveOperationHistory(&opHistory)
			if err != nil {
				log.Fatal(err)
			}

			// fmt.Printf("发送者：%s\n", sender)

			// //查询所有输出用户及金额
			// receiptList, amountList, err := getAllReceiptAddressAndAmount(tx)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// receiptAndAmountList := make([]string, len(receiptList))
			// for i := 0; i < len(receiptList); i++ {
			// 	receiptAndAmountList[i] = fmt.Sprintf("[%s:%v]", receiptList[i], amountList[i])
			// }

			// fmt.Printf("发送者：%s\n", sender)

			// // 解析交易信息
			// fmt.Printf("交易ID：%s\n", txID)
			// fmt.Printf("交易哈希：%s\n", tx.Hash)
			// fmt.Printf("接收者：%s\n", strings.Join(receiptAndAmountList, ";")) //接收者有多个，因此这里先不解析
			// fmt.Printf("交易时间：%s\n", time.Unix(tx.Time, 0))
			// getInscription(tx)
			// fmt.Println("--------------------------")
		}

		currentHeight++
		if err := models.SaveChainHeight(oldHight, currentHeight); err != nil {
			zlog.Zlog.Warn("update chain height err height", zap.Int64("currentHeight", currentHeight),
				zap.Error(err))
		}

	}
}

// 获取交易的发送者
func getSenderAddress(verboseTx *btcjson.TxRawResult) (string, error) {
	// 遍历输入,随便找到一个转入的交易
	for _, vin := range verboseTx.Vin {
		if len(vin.Txid) == 0 {
			continue
		}
		// 获取PrevOut的交易ID和输出索引
		prevTxID := vin.Txid
		prevOutIndex := vin.Vout

		//根据交易的Id查询交易详情
		prevTxHash, err := chainhash.NewHashFromStr(prevTxID)
		if err != nil {
			log.Fatal(err)
			return "", err

		}
		preTxVerbose, err := node.BtcClient.GetClient().GetRawTransactionVerbose(prevTxHash)
		if err != nil {
			log.Fatal(err)
			return "", err

		}

		// 根据交易和交易输出的索引 查找接收者和金额
		receipt, _, err := getReceiptAddressAndAmount(preTxVerbose, int32(prevOutIndex))
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		return receipt, nil
	}
	return "", nil
}

// 获取交易的所有接收者和金额
func getAllReceiptAddressAndAmount(verboseTx *btcjson.TxRawResult) (receiptList []string, amountList []float64, err error) {
	// 初始化切片
	receiptList = make([]string, len(verboseTx.Vout))
	amountList = make([]float64, len(verboseTx.Vout))

	// 遍历所有输出
	for index, _ := range verboseTx.Vout {
		// 根据交易和交易输出的索引 查找接收者和金额
		receipt, amount, err := getReceiptAddressAndAmount(verboseTx, int32(index))
		if err != nil {
			return nil, nil, err
		}

		receiptList[index] = receipt
		amountList[index] = amount

	}

	return receiptList, amountList, nil
}

// 根据交易ID和交易输出的索引 查找接收者和金额
func getReceiptAddressAndAmount(txVerbose *btcjson.TxRawResult, outIndex int32) (receipt string, amount float64, err error) {
	//输出索引必须于输出交易的个数  TODO 这里可能报错
	if int(outIndex) >= len(txVerbose.Vout) {
		return "", 0, errors.New("outIndex to large")
	}
	txOut := txVerbose.Vout[outIndex]
	// 解析ScriptPubKey
	scriptPubKey, err := hex.DecodeString(txOut.ScriptPubKey.Hex)
	if err != nil {
		return "", 0, err
	}

	// 解析地址
	_, addrs, _, err := txscript.ExtractPkScriptAddrs(scriptPubKey, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
		return "", 0, err
	}
	if len(addrs) == 0 {
		fmt.Printf("  %s\n", addrs)
		return "A", 0, nil //TODO
	}
	return addrs[0].String(), txOut.Value, nil
}

// 获取铭文
func getInscription(tx *btcjson.TxRawResult) (inscription string, err error) {
	// 检查输入脚本是否为空
	if len(tx.Vin) == 0 || len(tx.Vin[0].Witness) == 0 {
		fmt.Println("输入脚本为空")
		return
	}

	// // 初始化切片，避免重新分配内存
	// witnessData := make([][]byte, len(tx.Vin[0].Witness))
	// // 将输入脚本中的witness数据转换为[][]byte类型
	// for i, witness := range tx.Vin[0].Witness {
	// 	witnessData[i] = []byte(witness)
	// }

	// witnessDataMerged := bytes.Join(witnessData, []byte{})

	// decodedDatarsp, err := hex.DecodeString(string(witnessDataMerged))
	// if err != nil {
	// 	return "", err
	// }
	// fmt.Println("decodedDatarsp", decodedDatarsp)

	// 将输入脚本中的witness数据转换为[][]byte类型
	for _, witness := range tx.Vin[0].Witness {
		decodedData, err := hex.DecodeString(witness)
		if err != nil {
			return "", err
		}
		witnessStr := string(decodedData)
		if !strings.Contains(witnessStr, "ord") || !strings.Contains(witnessStr, "op") {
			continue
		}
		witnessFinal := witnessStr[strings.LastIndex(witnessStr, "{") : strings.LastIndex(witnessStr, "}")+1]
		fmt.Println("铭文内容：", witnessFinal)
		return witnessFinal, nil

	}
	return "", nil

}
