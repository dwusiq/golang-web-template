package models

import (
	"errors"
	"time"
)

type (
	BasicModel struct {
		Status  int32     `json:"status"`
		Deleted bool      `json:"-"`
		Created time.Time `orm:"auto_now_add; type(datetime)" json:"-"`
		Updated time.Time `orm:"auto_now; type(datetime)" json:"-"`
	}
)

// 错误定义
// example:
// var BalanceErr = errors.New("余额不足")
var (
	DateTooLong = errors.New("时间跨度太长")
)

var (
	BalanceOutError      = errors.New("账户余额不足")
	NewsNotExistError    = errors.New("内容不存在或已被删除")
	PriceChangedError    = errors.New("当前价格已被改变,请刷新页面后重试")
	TimePackageOutError  = errors.New("托管包不足,请购买后使用")
	LimitAwardExpired    = errors.New("奖励金已失效")
	ProductNotExist      = errors.New("商品不存在或还未上市")
	PictureNotExistError = errors.New("图片不存在或已被删除")
	NodeIsExistError     = errors.New("矿机已经添加")
)
