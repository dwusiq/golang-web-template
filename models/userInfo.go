package models

import (
	"strings"

	"github.com/astaxie/beego/orm"
)
//用户表
type UserInfo struct {
	BasicModel
	Id      int64  `orm:"pk;auto" json:"id"`
	Address string `json:"address"`
}

/**
 * @dev 根据钱包地址查询,如果查不到就写入
 * @param address 钱包地址
 */
func FindOrSaveAddress(address string) (userInfo *UserInfo, err error) {
	o := orm.NewOrm()
	//查询
	queryRsp := &UserInfo{}
	err = o.QueryTable("UserInfo").Filter("address", address).One(queryRsp)
	if queryRsp.Address != "" {
		return queryRsp, err
	}

	if err != nil && !strings.Contains(err.Error(), "no row found") {
		return queryRsp, err
	}

	//写入
	_, err = o.Insert(&UserInfo{Address: address})
	if err != nil {
		return nil, err
	}

	//查询结果返回
	return FindByAddress(address)
}

/**
 * @dev 根据钱包地址查询
 * @param address 钱包地址
 */
func FindByAddress(addr string) (userInfo *UserInfo, err error) {
	o := orm.NewOrm()
	queryRsp := &UserInfo{}
	err = o.QueryTable("UserInfo").Filter("address", addr).One(queryRsp)
	return queryRsp, err
}

/**
 * @dev 根据ID查询
 * @param id 用户ID
 */
func FindUserById(id int64) (userInfo *UserInfo, err error) {
	o := orm.NewOrm()

	//条件
	cond := orm.NewCondition()
	cond = cond.And("id", id)

	//构建查询对象
	qs := o.QueryTable("UserInfo")
	qs.SetCond(cond)

	//查询
	err = qs.One(userInfo)
	return
}
