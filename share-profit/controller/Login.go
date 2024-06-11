package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-acme/lego/v3/log"
	"net/http"
	"share-profit/conf"
	"share-profit/dto"
	"share-profit/entity"
	"share-profit/models"
	"share-profit/pkg/googleAuthenticator"
	"share-profit/result"
	"share-profit/utils"
	"strconv"
	"time"
)

type Login struct {
	*models.LoginModels
}
func NewLogin(login models.LoginModels) (*Login){
	return &Login{
		&login,
	}
}
//@description 用户登录
//@accept json
//@Param loginDto body dto.LoginDto true "loginDto"
//@Success 200 {object} gin.H
//@router /login [post]
func (login *Login) DoLogin(ctx *gin.Context) {

	var loginDto dto.LoginDto

	if err := ctx.ShouldBindJSON(&loginDto); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "data": nil, "msg": err.Error()})
		ctx.Abort()
		return

	}
	sysUser := new(entity.SysUser)
	if ib, err := login.Engine.Where("sys_user.username=?", loginDto.Username).Get(sysUser); err != nil || ib == false {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "data": nil, "msg": "用户名错误"})
		ctx.Abort()
		return
	}
	//密码b不匹配
	if utils.EncodePassword(loginDto.Username, loginDto.Password) != sysUser.Password {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "data": nil, "msg": "密码错误"})
		ctx.Abort()
		return
	}

	//1:平台管理员(sysUser.DeptId == 0 && sysUser.SuperAdmin==1)
	//2:商户系统用户(管理员<sysUser.DeptId != 0 && sysUser.SuperAdmin==1>，非管理员<sysUser.DeptId != 0 && sysUser.SuperAdmin==0>)
	//3:普通用户(sysUser.DeptId == 0 && sysUser.SuperAdmin==0)

	//账号密码输入成功之后（商户管理员账号被审核未通过无法登录）
	if sysUser.DeptId != 0 && sysUser.SuperAdmin==1 && sysUser.Verified==0 {

		ctx.JSON(200,result.Fail(errors.New("商户审核未通过")))
		ctx.Abort()
		return
	}

	//签名生成token,设置有效期为一周
	token, err := utils.Sign(loginDto.Username, conf.TOKEN_EFFECT_TIME)
	if err != nil {
		ctx.String(500, err.Error())
		ctx.Abort()
		return
	}
	//用户基本信息缓存到redis,并设置有效期为token有效期的2倍
	userInfo, _ := json.Marshal(sysUser)
	err = login.SetNX(conf.TOKEN_PREFIX+token, string(userInfo), 2*time.Second*time.Duration(conf.TOKEN_EFFECT_TIME)).Err()
	if err != nil {
		ctx.String(500, err.Error())
		ctx.Abort()
		return
	}
	//获取是否已绑定私钥
	sectet :=  login.GetSysUserVerification(strconv.FormatInt(sysUser.Id,10))
	if sectet==""{
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": token, "msg": "未绑定"})
		return
	}

	//加密生成token，并存到redis中
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": token, "msg": "ok"})
}

//@description 用户户注册
//@accept json
//@Param sysUser body entity.SysUser true "sysUser"
//@Success 200 {object} gin.H
//@router /register [post]
func (login *Login) DoRegister(ctx *gin.Context) {
	var sysUser entity.SysUser
	if err := ctx.ShouldBindJSON(&sysUser); err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
		ctx.Abort()
		return
	}
	//用户名和密码md5加密
	sysUser.Id = utils.Node().Generate().Int64()
	sysUser.Password = utils.EncodePassword(sysUser.Username, sysUser.Password)
	sysUser.Creator = 1067246875800000001
	sysUser.Updater = 1067246875800000001
	sysUser.CreateDate = time.Now()
	sysUser.UpdateDate = time.Now()
	sysUser.AdminId = utils.Node().Generate().Int64()
	//注册lotus钱包
	//addr, err := login.Rpc.WalletNew(context.TODO(), types.KTBLS)
	//if err != nil {
	//	ctx.JSON(http.StatusOK, result.Fail(err))
	//	ctx.Abort()
	//	return
	//}
	//sysUser.Wallet=addr.String()
	if _, err := login.Engine.Insert(&sysUser); err != nil {
		ctx.JSON(http.StatusOK, result.Fail(err))
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, result.Ok(nil))
}

//@description 商户注册申请
//@accept json
//@Param MerchantsDto body dto.MerchantsDto true "MerchantsDto"
//@Success 200 {object} gin.H
//@router /merchants/register [post]
func (login *Login) MerchantsRegister(ctx *gin.Context) {

	var merchantsDto dto.MerchantsDto

	if err := ctx.ShouldBindJSON(&merchantsDto); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "data": nil, "msg": err.Error()})
		return
	}

	log.Printf("申请对象:%v", merchantsDto)

	//创建商户组织
	dept := entity.SysDept{
		Id:         utils.Node().Generate().Int64(),
		Name:       merchantsDto.Name,
		CreateDate: time.Now(),
		UpdateDate: time.Now(),
	}
	//商户组织关联管理员账号
	sysUser := entity.SysUser{
		Id:       utils.Node().Generate().Int64(),
		Username: merchantsDto.Username,
		Password: utils.EncodePassword(merchantsDto.Username, merchantsDto.Password),
		RealName: merchantsDto.RealName,
		HeadUrl:  merchantsDto.HeadUrl,
		Gender:   merchantsDto.Gender,
		Email:    merchantsDto.Email,
		Mobile:   merchantsDto.Mobile,
		DeptId:   dept.Id,
		//组织机构超级管理员:0否1是
		SuperAdmin: 1,
		Remark:     merchantsDto.Remark,
		CreateDate: time.Now(),
		UpdateDate: time.Now(),
	}

	//开启事务，先创建组织 再创建商户账号
	tx := login.Engine.NewSession()

	defer tx.Close()

	// add Begin() before any action
	if err := tx.Begin(); err != nil {
		return
	}

	log.Println("---------开启事务-----------")

	if _, err := tx.Insert(&dept); err != nil {
		tx.Rollback()
		log.Printf("提交申请人错误：", err)
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	if _, err := tx.Insert(&sysUser); err != nil {
		tx.Rollback()
		log.Printf("提交组织错误：", err)
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	err := tx.Commit()

	if err != nil {
		log.Printf("提交错误：", err)
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": merchantsDto, "msg": "ok"})
}

func (login *Login) CreateSecret(ctx *gin.Context) {
	//判断用户是否已经绑定谷歌验证吗
	sysu,_ := utils.GetSubject(login.Client,ctx)

	var res result.Result
	var goo models.GoogleAuth

	if len(sysu.Verification) <= 0 {
		ga := googleAuthenticator.NewGAuth()
		secret, _ := ga.CreateSecret(16)

		goo.GoogleSecret = secret
		goo.GoogleAccount = "yungo:" + sysu.Username
		goo.URL = "otpauth://totp/" + goo.GoogleAccount + "?secret=" + secret
		goo.IsBind = false
		res.Data = goo
		res.Msg = "获取成功！"
	} else {
		goo.IsBind = true
		res.Data = goo
		res.Msg = "已经绑定！"
	}

	ctx.JSON(http.StatusOK, res)
}

func (login *Login) BindCode(ctx *gin.Context) {
	//var data models.GoogleAuth
	GoogleSecret := ctx.Request.FormValue("googleSecret")
	code := ctx.Request.FormValue("code")
	log.Println("BindCode VerifyCode:", GoogleSecret,code)

	ga := googleAuthenticator.NewGAuth()
	isOK, _ := ga.VerifyCode(GoogleSecret, code, 1)
	if isOK {
		sysu,_ := utils.GetSubject(login.Client,ctx)
		param := make(map[string]string)
		param["verification"] = GoogleSecret

		param["id"] = strconv.FormatInt(sysu.Id,10)
		err := login.UpdateSysUser(param)
		if err != nil{
			ctx.JSON(http.StatusOK, result.Fail(errors.New("验证失败！")))
			return
		}

	} else {

		ctx.JSON(http.StatusOK, result.Fail(errors.New("验证码无效！")))
		return
	}
	tokenStr:= ctx.GetHeader(conf.API_KEY)
	user,err := utils.GetSubjectByTokenStr(tokenStr,login.Client)
	if err!=nil{
		ctx.JSON(http.StatusOK, result.Fail(errors.New("验证失败！")))
		return
	}
	user.Verified = 1

	userInfo, _ := json.Marshal(user)
	err = login.Set(conf.TOKEN_PREFIX+tokenStr, string(userInfo), 2*time.Second*time.Duration(conf.TOKEN_EFFECT_TIME)).Err()
	if err != nil {
		ctx.String(500, err.Error())
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, result.Ok("绑定成功！"))
}

//验证谷歌验证码
func (login *Login)VerifyCode(ctx *gin.Context) {

	code := ctx.Request.FormValue("code")
	tokenStr:= ctx.GetHeader(conf.API_KEY)
	user,err := utils.GetSubjectByTokenStr(tokenStr,login.Client)
	if err!=nil{
		ctx.JSON(http.StatusOK, result.Fail(errors.New("验证失败！")))
		return
	}
	sectet :=  login.GetSysUserVerification(strconv.FormatInt(user.Id,10))
	ga := googleAuthenticator.NewGAuth()
	isOk ,_ := ga.VerifyCode(sectet,code,1)
	if !isOk {
		ctx.JSON(http.StatusOK, result.Fail(errors.New("验证失败！")))
		return
	}
	user.Verified = 1

	userInfo, _ := json.Marshal(user)
	err = login.Set(conf.TOKEN_PREFIX+tokenStr, string(userInfo), 2*time.Second*time.Duration(conf.TOKEN_EFFECT_TIME)).Err()
	if err != nil {
		ctx.String(500, err.Error())
		ctx.Abort()
		return
	}



	ctx.JSON(http.StatusOK, result.Ok("验证成功！"))
}
