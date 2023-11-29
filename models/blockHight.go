package models

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego/orm"
)

type BlockHight struct {
	Id    int   `orm:"pk;auto" json:"id"`
	Hight int64 `json:"hight"`
}

func ReadChainHeight() (int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("blockHight")

	var blockHight BlockHight
	// err := qs.One(&blockHight)
	qs.One(&blockHight)

	// if err != nil {
	// 	return 0, err
	// }
	return strconv.ParseInt(fmt.Sprint(blockHight.Hight), 10, 64)

	// return nil

	// if !Has(key) {
	// 	return 0, nil
	// }
	// height, err := database.Leveldb.Get([]byte(key), nil)
	// if err != nil {
	// 	return 0, err
	// }
	// return strconv.ParseInt(string(height), 10, 64)
}

func SaveChainHeight(oldHight int64, height int64) error {
	o := orm.NewOrm()

	if oldHight == 0 {
		blockHight := BlockHight{Hight: height}
		_, err := o.Insert(&blockHight)
		return err

	} else {
		_, err := o.QueryTable("blockHight").Filter("Hight", oldHight).Update(orm.Params{
			"hight": height,
		})
		return err
	}

	// h := strconv.FormatInt(height, 10)
	// return database.Leveldb.Put([]byte(key), []byte(h), nil)
}

// func Has(key string) bool {
// 	blockHight, _ := ReadChainHeight(key)
// 	return blockHight > 0
// 	// has, _ := database.Leveldb.Has([]byte(key), nil)
// 	// return has
// }
