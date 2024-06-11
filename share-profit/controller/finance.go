package controller

import (
	"errors"
	"fmt"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"share-profit/conf"
	"share-profit/result"
	"share-profit/service"
	"share-profit/utils"
	"strconv"
	"strings"
)

type Finance struct {
	*service.FinanceService
}

func NewFinance(s service.FinanceService) *Finance {

	return &Finance{
		&s,
	}
}

func (fi *Finance) Upload(ctx *gin.Context) {

	file, err := ctx.FormFile("file")
	if err != nil {

	}
	tokenStr := ctx.GetHeader(conf.API_KEY)
	user, _ := utils.GetSubjectByTokenStr(tokenStr, fi.Client)
	f, _ := file.Open()

	if err = fi.UploadApply(f, file.Size, user); err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
	}
	ctx.JSON(http.StatusOK, result.Ok("导入成功"))
}
func (fi *Finance) ApplyList(ctx *gin.Context) {
	param := make(map[string]string)
	param["status"] = ctx.Request.FormValue("status")
	param["pageNum"] = ctx.DefaultQuery("pageNum", "1")
	param["pageSize"] = ctx.DefaultQuery("pageSize", "10")
	param["num"] = ctx.Request.FormValue("num")
	tokenStr := ctx.GetHeader(conf.API_KEY)
	user, _ := utils.GetSubjectByTokenStr(tokenStr, fi.Client)
	param["admin_id"] = strconv.FormatInt(user.AdminId, 10)

	re, err := fi.GetApplyList(param)
	res := utils.NewPageData(param, re)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
	}

	ctx.JSON(http.StatusOK, result.Ok(res))
}

//手动输入列表查询
func (fi *Finance) ManualList(ctx *gin.Context) {
	param := make(map[string]string)
	param["status"] = ctx.Request.FormValue("status")
	param["pageNum"] = ctx.DefaultQuery("pageNum", "1")
	param["pageSize"] = ctx.DefaultQuery("pageSize", "10")
	tokenStr := ctx.GetHeader(conf.API_KEY)
	user, _ := utils.GetSubjectByTokenStr(tokenStr, fi.Client)
	param["admin_id"] = strconv.FormatInt(user.AdminId, 10)

	re, err := fi.GetManualList(param)
	res := utils.NewPageData(param, re)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
	}

	ctx.JSON(http.StatusOK, result.Ok(res))
}

//批量手动列表转账
func (fi *Finance) SendManual(ctx *gin.Context) {
	fil := ctx.Request.FormValue("fil")
	to := ctx.Request.FormValue("to")
	id := ctx.Request.FormValue("id")
	tokenStr := ctx.GetHeader(conf.API_KEY)
	tos := strings.Split(to, ",")
	ids := strings.Split(id, ",")
	fils := strings.Split(fil, ",")
	var err error
	var message *types.SignedMessage
	user, _ := utils.GetSubjectByTokenStr(tokenStr, fi.Client)
	param := make(map[string]string)
	config := utils.GetConfig(fi.Engine, "charge", user)
	charge, err := strconv.ParseFloat(config, 64)
	if err != nil {
		charge = 0
	}
	for i, val := range tos {
		amout, _ := strconv.ParseFloat(fils[i], 64)
		fmt.Println("fil:", amout, ",to:", val, ",id:", ids[i])
		message, err = fi.Sends(val, amout-amout*charge, user)
		if err != nil {
			param["cid"] = "钱包地址不正确或者余额不足："
			param["error"] = err.Error()
			param["status"] = "0"
			param["value"] = "0"
		} else {
			param["status"] = "1"
			param["cid"] = message.Cid().String()
			param["error"] = ""
			param["value"] = strconv.FormatFloat(amout-amout*charge, 'f', -1, 64)
		}
		param["id"] = ids[i]
		param["to"] = val
		err = fi.UpdateManualInfo(param)
		if err != nil {
			break
		}
	}
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, result.Ok("转账成功！"))
}

//批量手动列表审核
func (fi *Finance) UpdateManual(ctx *gin.Context) {
	ids := ctx.Request.FormValue("ids")
	status := ctx.Request.FormValue("status")
	err := fi.UpdateManualList(ids, status)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, result.Ok("审核成功！"))
}
func (fi *Finance) ManualAdd(ctx *gin.Context) {
	param := make(map[string]string)
	param["user_name"] = ctx.Request.FormValue("user_name")
	param["amount"] = ctx.Request.FormValue("amount")
	param["address"] = ctx.Request.FormValue("address")
	tokenStr := ctx.GetHeader(conf.API_KEY)
	user, _ := utils.GetSubjectByTokenStr(tokenStr, fi.Client)
	param["admin_id"] = strconv.FormatInt(user.AdminId, 10)
	err := fi.ManualAdds(param)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, result.Ok("添加成功"))
}

//批量转账
func (fi *Finance) WalletSends(ctx *gin.Context) {
	fil := ctx.Request.FormValue("fil")
	to := ctx.Request.FormValue("to")
	id := ctx.Request.FormValue("id")
	numstr := ctx.Request.FormValue("num")

	tokenStr := ctx.GetHeader(conf.API_KEY)
	tos := strings.Split(to, ",")
	ids := strings.Split(id, ",")
	fils := strings.Split(fil, ",")
	var err error
	var message *types.SignedMessage
	user, _ := utils.GetSubjectByTokenStr(tokenStr, fi.Client)
	param := make(map[string]string)
	config := utils.GetConfig(fi.Engine, "charge", user)
	charge, err := strconv.ParseFloat(config, 64)
	if err != nil {
		charge = 0
	}
	//判断批次是否是在转账
	err, flag := fi.GetSendsFlag(tokenStr, numstr)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
		ctx.Abort()
		return
	}
	if !flag {
		ctx.JSON(http.StatusOK, result.Fail(errors.New("当前批次正在转账，请不要重复提交！")))
		ctx.Abort()
		return
	}

	num := 0
	for i, val := range tos {
		amout, _ := strconv.ParseFloat(fils[i], 64)
		fmt.Println("fil:", amout, ",to:", val, ",id:", ids[i])

		message, err = fi.Sends(val, amout-amout*charge, user)
		if err != nil {
			param["cid"] = "钱包地址不正确或者余额不足："
			param["error"] = err.Error()
			param["status"] = "0"
			num++
			param["value"] = "0"
		} else {
			param["status"] = "1"
			param["cid"] = message.Cid().String()
			param["error"] = ""
			param["value"] = strconv.FormatFloat(amout-amout*charge, 'f', -1, 64)
		}
		param["id"] = ids[i]
		param["to"] = val
		err = fi.UpdateTransferInfo(param)
		if err != nil {
			break
		}
	}
	//转账完毕 解锁批次
	if err := fi.SetSendsByRedis(tokenStr, numstr); err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
		ctx.Abort()
		return
	}
	if num != 0 {
		ctx.JSON(http.StatusOK, result.Ok(fmt.Sprintf("转账成功，%d个钱包异常", num)))
		return
	}
	ctx.JSON(http.StatusOK, result.Ok("转账成功！"))

}

//批量审核
func (fi *Finance) UpdateStatus(ctx *gin.Context) {
	ids := ctx.Request.FormValue("ids")
	status := ctx.Request.FormValue("status")
	err := fi.UpdateApplyList(ids, status)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, result.Ok("审核成功！"))
}

//获取钱包及余额
func (fi *Finance) WalletAndBalance(ctx *gin.Context) {

	tokenStr := ctx.GetHeader(conf.API_KEY)
	user, err := utils.GetSubjectByTokenStr(tokenStr, fi.Client)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
		ctx.Abort()
		return
	}
	addr, bigInt, err := fi.WalletBalance(strconv.FormatInt(user.AdminId, 10))
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
		ctx.Abort()
		return
	}
	bls := utils.NanoOrAttoToFIL(bigInt.String(), utils.AttoFIL)
	res := make(map[string]interface{})
	res["address"] = addr
	res["balance"] = bls
	ctx.JSON(http.StatusOK, result.Ok(res))
}

//获取批次列表
func (fi *Finance) BatchList(ctx *gin.Context) {
	param := make(map[string]string)
	param["pageNum"] = ctx.DefaultQuery("pageNum", "1")
	param["pageSize"] = ctx.DefaultQuery("pageSize", "10")
	param["status"] = ctx.Request.FormValue("status")

	tokenStr := ctx.GetHeader(conf.API_KEY)
	user, _ := utils.GetSubjectByTokenStr(tokenStr, fi.Client)
	param["admin_id"] = strconv.FormatInt(user.AdminId, 10)

	re, err := fi.GetBatchList(param)
	res := utils.NewPageData(param, re)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
	}

	ctx.JSON(http.StatusOK, result.Ok(res))
}

//按批次转账
func (fi *Finance) BatchSends(ctx *gin.Context) {
	id := ctx.Request.FormValue("id")

	tokenStr := ctx.GetHeader(conf.API_KEY)
	user, _ := utils.GetSubjectByTokenStr(tokenStr, fi.Client)

	//判断批次是否是在转账
	err, flag := fi.GetSendsFlag(tokenStr, id)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
		ctx.Abort()
		return
	}
	if !flag {
		ctx.JSON(http.StatusOK, result.Fail(errors.New("当前批次正在转账，请不要重复提交！")))
		ctx.Abort()
		return
	}
	num, err := fi.BatchSendsList(id, user)
	//转账完毕 解锁批次
	if err := fi.SetSendsByRedis(tokenStr, id); err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
		ctx.Abort()
		return
	}
	if num != 0 {
		ctx.JSON(http.StatusOK, result.Ok(fmt.Sprintf("转账成功，%d个钱包异常", num)))
		return
	}
	ctx.JSON(http.StatusOK, result.Ok("转账成功！"))

}

//按批次拒绝
func (fi *Finance) BatchRefuse(ctx *gin.Context) {
	ids := ctx.Request.FormValue("ids")
	status := ctx.Request.FormValue("status")
	err := fi.BatchRefuseList(ids, status)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, result.Ok("审核成功！"))
}

//获取配置列表
func (fi *Finance) GetConfig(ctx *gin.Context) {

	tokenStr := ctx.GetHeader(conf.API_KEY)
	user, _ := utils.GetSubjectByTokenStr(tokenStr, fi.Client)
	res := fi.GetConfigList(strconv.FormatInt(user.AdminId, 10))

	ctx.JSON(http.StatusOK, result.Ok(res))
}

//设置配置列表
func (fi *Finance) SetConfig(ctx *gin.Context) {

	tokenStr := ctx.GetHeader(conf.API_KEY)
	user, _ := utils.GetSubjectByTokenStr(tokenStr, fi.Client)
	param := make(map[string]string)
	param["charge"] = ctx.Request.FormValue("charge")
	res := fi.SetConfigList(strconv.FormatInt(user.AdminId, 10), param)

	ctx.JSON(http.StatusOK, result.Ok(res))
}
