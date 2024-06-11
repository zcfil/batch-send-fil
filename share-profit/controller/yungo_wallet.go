package controller

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/filecoin-project/go-address"
	"github.com/gin-gonic/gin"
	"net/http"
	"share-profit/conf"
	"share-profit/pkg"
	"share-profit/result"
	"share-profit/service"
	"share-profit/utils"
	"strconv"
)

type SendParams struct {
	To string `json:"to" binding:"required"`
	Fil float64 `json:"fil" binding:"required"`
}

type Wallet struct{
	*service.YungoService
}

func NewWallet(ys service.YungoService) *Wallet  {
	return &Wallet{
		&ys,
	}
}


//创建助记词
func (y *Wallet) GetMnemonic(ctx *gin.Context){

	res := pkg.CreateMnemonic()

	ctx.JSON(http.StatusOK,result.Ok(utils.AesEncrypt([]byte(res))))

}
//创建钱包地址 pkg.CreateMnemonic()
func (y *Wallet) NewWallet(ctx *gin.Context){

	tokenStr := ctx.GetHeader(conf.API_KEY)
	user, err := utils.GetSubjectByTokenStr(tokenStr, y.Client)
	if err != nil {
		ctx.JSON(http.StatusOK,result.Fail(err))
		ctx.Abort()
		return
	}

	param := make(map[string]string)
	param["admin_id"] = strconv.FormatInt(user.AdminId,10)
	//助记词
	str := ctx.Request.FormValue("mnemonic")
	mne,_ := hex.DecodeString(str)
	param["mnemonic"] = string(utils.AesDecrypt(mne))
	param["mnemonic_ase"] = str
	addr, err := y.YungoService.WalletNew(param)
	if err != nil {
		ctx.JSON(http.StatusOK,result.Fail(err))
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK,result.Ok(addr))

}

//获取钱包列表
func (y *Wallet) WalletList(ctx *gin.Context){
	tokenStr := ctx.GetHeader(conf.API_KEY)
	user, err := utils.GetSubjectByTokenStr(tokenStr, y.Client)
	if err != nil {
		ctx.JSON(http.StatusOK,result.Fail(err))
		ctx.Abort()
		return
	}
	param := make(map[string]string)
	param["pageNum"] = ctx.DefaultQuery("pageNum","1")
	param["pageSize"] = ctx.DefaultQuery("pageSize","10")
	param["admin_id"] = strconv.FormatInt(user.AdminId,10)
	re, err := y.GetWalletList(param)
	if err!=nil{
		ctx.JSON(http.StatusOK, result.Fail(err))
		return
	}
	res := utils.NewPageData(param,re)
	ctx.JSON(http.StatusOK, result.Ok(res))
}
//删除钱包
func (y *Wallet)DelWallet(ctx *gin.Context){
	tokenStr := ctx.GetHeader(conf.API_KEY)
	user, err := utils.GetSubjectByTokenStr(tokenStr, y.Client)
	if err != nil {
		ctx.JSON(http.StatusOK,result.Fail(err))
		ctx.Abort()
		return
	}
	param := make(map[string]string)
	param["wallet"] =  ctx.Request.FormValue("wallet")
	param["admin_id"] = strconv.FormatInt(user.AdminId,10)
	if y.IsPrivateKey(param){
		addr,err := address.NewFromString(param["wallet"])
		if err != nil {
			ctx.JSON(http.StatusOK,result.Fail(err))
			ctx.Abort()
			return
		}
		//如果没人拥有，就删除
		if !y.IsUseWallet(param){
			err = y.YungoRpc.Rpc.WalletDelete(ctx,addr )
			if err!=nil{
				ctx.JSON(http.StatusOK, result.Fail(err))
				return
			}
		}
		err = y.DeleteWallet(param)
		if err!=nil{
			ctx.JSON(http.StatusOK, result.Fail(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, result.Ok("删除成功！"))
}
//设置默认地址
func (y *Wallet)SetWallet(ctx *gin.Context){
	tokenStr := ctx.GetHeader(conf.API_KEY)
	user, err := utils.GetSubjectByTokenStr(tokenStr, y.Client)
	if err != nil {
		ctx.JSON(http.StatusOK,result.Fail(err))
		ctx.Abort()
		return
	}
	param := make(map[string]string)
	param["wallet"] =  ctx.Request.FormValue("wallet")
	param["admin_id"] = strconv.FormatInt(user.AdminId,10)
	err = y.SetDefault(param)
	if err!=nil{
		ctx.JSON(http.StatusOK, result.Fail(err))
		return
	}

	ctx.JSON(http.StatusOK, result.Ok("删除成功！"))
}

//通过私钥导入
func (y *Wallet) ImportWallet(ctx *gin.Context){

	tokenStr := ctx.GetHeader(conf.API_KEY)
	user, err := utils.GetSubjectByTokenStr(tokenStr, y.Client)
	if err != nil {
		ctx.JSON(http.StatusOK,result.Fail(err))
		ctx.Abort()
		return
	}

	param := make(map[string]string)
	param["admin_id"] = strconv.FormatInt(user.AdminId,10)
	param["private_key"] = ctx.Request.FormValue("private_key")


	addr, err := y.YungoService.ImprotWalletAddress(param)
	if err != nil {
		ctx.JSON(http.StatusOK,result.Fail(err))
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK,result.Ok(addr))

}

//导出私钥
func (y *Wallet) ExportWallet(ctx *gin.Context){

	tokenStr := ctx.GetHeader(conf.API_KEY)
	user, err := utils.GetSubjectByTokenStr(tokenStr, y.Client)
	if err != nil {
		ctx.JSON(http.StatusOK,result.Fail(err))
		ctx.Abort()
		return
	}

	param := make(map[string]string)
	param["admin_id"] = strconv.FormatInt(user.AdminId,10)
	param["wallet"] = ctx.Request.FormValue("wallet")
	if !y.IsPrivateKey(param) {
		ctx.JSON(http.StatusOK,result.Fail(errors.New("该钱包私钥不存在！")))
		ctx.Abort()
		return
	}
	if utils.CheckPrefix(param["wallet"],"f3"){
		var addr address.Address
		addr.Scan(param["wallet"])
		keyinfo,_ := y.Rpc.WalletExport(context.TODO(),addr)
		keystr,_ := json.Marshal(*keyinfo)
		if err!=nil{
			ctx.JSON(http.StatusOK,result.Fail(err))
			ctx.Abort()
			return
		}
		//加密
		ctx.JSON(http.StatusOK,result.Ok(utils.AesEncrypt([]byte(hex.EncodeToString(keystr)))))
		return

	}
	keystr,err := pkg.ExportKey(param["wallet"])
	if err!=nil{
		ctx.JSON(http.StatusOK,result.Fail(err))
		ctx.Abort()
		return
	}
	//加密
	ctx.JSON(http.StatusOK,result.Ok(utils.AesEncrypt([]byte(hex.EncodeToString(keystr)))))


}

//导出助记词
func (y *Wallet) ExportMnemonic(ctx *gin.Context){

	tokenStr := ctx.GetHeader(conf.API_KEY)
	user, err := utils.GetSubjectByTokenStr(tokenStr, y.Client)
	if err != nil {
		ctx.JSON(http.StatusOK,result.Fail(err))
		ctx.Abort()
		return
	}

	param := make(map[string]string)
	param["admin_id"] = strconv.FormatInt(user.AdminId,10)
	param["wallet"] = ctx.Request.FormValue("wallet")
	if !y.IsPrivateKey(param){
		ctx.JSON(http.StatusOK,result.Fail(errors.New("该钱包私钥不存在！")))
		ctx.Abort()
		return
	}

	addr, err := y.YungoService.ExportMnemonicByAddres(param)
	if err != nil {
		ctx.JSON(http.StatusOK,result.Fail(err))
		ctx.Abort()
		return
	}
	//加密
	//ctx.JSON(http.StatusOK,result.Ok(utils.AesEncrypt([]byte(addr))))
	res,_ := hex.DecodeString(addr)
	ctx.JSON(http.StatusOK,result.Ok(res))

}


//通过助记词导入
func (y *Wallet) ImportMnemonic(ctx *gin.Context){

	tokenStr := ctx.GetHeader(conf.API_KEY)
	user, err := utils.GetSubjectByTokenStr(tokenStr, y.Client)
	if err != nil {
		ctx.JSON(http.StatusOK,result.Fail(err))
		ctx.Abort()
		return
	}

	param := make(map[string]string)
	param["admin_id"] = strconv.FormatInt(user.AdminId,10)

	str := ctx.Request.FormValue("mnemonic")
	if err != nil {
		ctx.JSON(http.StatusOK,result.Fail(err))
		ctx.Abort()
		return
	}
	if y.IsMnemonic(str,param["admin_id"]){
		ctx.JSON(http.StatusOK,result.Fail(errors.New("该账号已导入！")))
		ctx.Abort()
		return
	}
	//解密
	mne,_ := hex.DecodeString(str)
	param["mnemonic"] = string(utils.AesDecrypt(mne))
	if !y.CheckMnemonic(param["mnemonic"]){
		ctx.JSON(http.StatusOK,result.Fail(errors.New("请输入正确的12位助记词！空格隔开")))
		ctx.Abort()
		return
	}
	param["mnemonic_ase"] = str

	addr, err := y.YungoService.ImprotMnemonicAddress(param)
	if err != nil {
		ctx.JSON(http.StatusOK,result.Fail(err))
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK,result.Ok(addr))

}
