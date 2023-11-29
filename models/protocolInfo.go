package models

import (
	"strings"

	"github.com/astaxie/beego/orm"
)

// 协议信息表
type ProtocolInfo struct {
	BasicModel
	Id            int64  `orm:"pk;auto" json:"id"`                  //自增ID
	Deployer      int64  `json:"deployer" binding:"required"`       //部署的用户id, address_info表
	DeployedTxId  string `json:"deployed_tx_id" binding:"required"` //部署的交易id
	DeployedBlock int64  `json:"deployed_block" binding:"required"` //部署的区块
	ProtocolType  int    `json:"protocol_type"  binding:"required"` //参考Protocol枚举
	TickName      string `json:"tick_name"  binding:"required"`     //票据名称，即代币或nft名称
	MintMax       int64  `json:"mint_max"`                          //总共最大发行量
	MintLimit     int64  `json:"mint_limit"`                        //每次最大铸造量
	// DeployedSequence int    `json:"deployed_sequence"`                 //部署次序，相同tick名字的前提下，只有第一个部署才有效，后续重名部署均视为无效
	ProtocolJson string `json:"protocol_json"  binding:"required"` //如 {"p":"brc-20","op":"deploy","tick":"ordi","max":"10000","max":"10"}
}

/**
 * @dev 查询,如果查不到就写入
 */
func SaveProtocol(protocolInfo *ProtocolInfo) (id int64, err error) {
	o := orm.NewOrm()
	queryRsp := &ProtocolInfo{}
	err = o.QueryTable("ProtocolInfo").Filter("TickName", protocolInfo.TickName).One(queryRsp)
	if queryRsp.TickName != "" {
		return queryRsp.Id, err
	}

	if err != nil && !strings.Contains(err.Error(), "no row found") {
		return 0, err
	}

	//写入
	id, err = o.Insert(protocolInfo)
	return id, err
}

/**
 * @dev 查询,如果查不到就写入
 */
func SaveProtocolInfo(protocolInfo *ProtocolInfo) error {
	_, err := FindProtocolInfoByTxId(protocolInfo.DeployedTxId)
	if err != nil {
		return err
	}

	//写入
	o := orm.NewOrm()
	_, err = o.Insert(protocolInfo)
	return err
}

/**
 * @dev 根据交易ID查询操作历史记录
 * @param txId 钱包地址
 */
func FindProtocolInfoByTxId(txId string) (protocolInfo *ProtocolInfo, err error) {
	o := orm.NewOrm()

	//条件
	cond := orm.NewCondition()
	cond = cond.And("deployedTxId", txId)

	//构建查询对象
	qs := o.QueryTable("ProtocolInfo")
	qs.SetCond(cond)

	//查询
	err = qs.One(protocolInfo)
	return
}

/**
 * @dev 根据ID查询
 * @param id 记录ID
 */
func FindProtocolInfoById(id int) (protocolInfo *ProtocolInfo, err error) {
	o := orm.NewOrm()

	//条件
	cond := orm.NewCondition()
	cond = cond.And("id", id)

	//构建查询对象
	qs := o.QueryTable("ProtocolInfo")
	qs.SetCond(cond)

	//查询
	err = qs.One(protocolInfo)
	return
}
