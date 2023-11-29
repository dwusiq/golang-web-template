package handler

import (
	"github.com/gin-gonic/gin"
	// "walletSynV2/handler/model"
	// "walletSynV2/models"
	// "walletSynV2/pkg"
)

func GetBlackList(c *gin.Context) {
	// list := models.GetBlackList()
	// resp := pkg.MakeResp(pkg.Success, list)
	// c.JSON(resp.HttpCode, resp)
	return
}

func CheckAccount(c *gin.Context) {
	// var ar model.CheckAccount
	// if err := c.ShouldBindJSON(&ar); err != nil {
	// 	resp := pkg.MakeResp(pkg.ParamsError, nil)
	// 	c.JSON(resp.HttpCode, resp)
	// 	return
	// }
	// account := models.CheckAccount(ar.Account)
	// resp := pkg.MakeResp(pkg.Success, account)
	// c.JSON(resp.HttpCode, resp)
	return
}

func DeleteAccount(c *gin.Context) {
	// var ar model.CheckAccount
	// if err := c.ShouldBindJSON(&ar); err != nil {
	// 	resp := pkg.MakeResp(pkg.ParamsError, nil)
	// 	c.JSON(resp.HttpCode, resp)
	// 	return
	// }

	// if err := models.DeleteAccount(ar.Account); err != nil {
	// 	resp := pkg.MakeResp(pkg.ParamsError, false)
	// 	c.JSON(resp.HttpCode, resp)
	// 	return
	// }
	// resp := pkg.MakeResp(pkg.Success, true)
	// c.JSON(resp.HttpCode, resp)
	return
}
