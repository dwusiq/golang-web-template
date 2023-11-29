package models

import (
	"strings"

	"github.com/astaxie/beego/orm"
)

// 操作的历史记录表
type OperationHistory struct {
	BasicModel
	Id            int64   `orm:"pk;auto" json:"id"`                  //记录ID
	ProtocolId    int64   ` json:"protocol_id"`                      //协议记录的ID
	OperationType int     ` json:"operation_type"`                   //操作类型，参考枚举
	OperatedTxId  string  `json:"operated_tx_id" binding:"required"` //操作交易的txId
	OperatedBlock int64   `json:"operated_block" binding:"required"` //操作的区块
	FromUser      int64   ` json:"from_user"`                        //发起者
	ToUser        int64   ` json:"to_user"`                          //接收者
	Amount        float64 ` json:"amount"`                           //操作数量
	// ConfirmedTime time.Time ` type(datetime)" json:"-"`                //区块确认的时间()  TODO 确认时间
}

/**
 * @dev 查询,如果查不到就写入
 * @param address 钱包地址
 */
func SaveOperationHistory(operationHistory *OperationHistory) error {
	o := orm.NewOrm()

	queryRsp := &OperationHistory{}
	err := o.QueryTable("OperationHistory").Filter("OperatedTxId", operationHistory.OperatedTxId).One(queryRsp)
	if queryRsp.OperatedTxId != "" {
		return err
	}

	if err != nil && !strings.Contains(err.Error(), "no row found") {
		return err
	}

	//写入
	_, err = o.Insert(operationHistory)
	return err
}

/**
 * @dev 根据交易ID查询操作历史记录
 * @param txId 钱包地址
 */
func FindOperationHistoryByTxId(txId string) (operationHistory *OperationHistory, err error) {
	o := orm.NewOrm()

	//条件
	cond := orm.NewCondition()
	cond = cond.And("operatedTxId", txId)

	//构建查询对象
	qs := o.QueryTable("OperationHistory")
	qs.SetCond(cond)

	//查询
	err = qs.One(operationHistory)
	return
}

/**
 * @dev 根据ID查询
 * @param id 记录ID
 */
func FindOperationHistoryById(id int) (operationHistory *OperationHistory, err error) {
	o := orm.NewOrm()

	//条件
	cond := orm.NewCondition()
	cond = cond.And("id", id)

	//构建查询对象
	qs := o.QueryTable("OperationHistory")
	qs.SetCond(cond)

	//查询
	err = qs.One(operationHistory)
	return
}
