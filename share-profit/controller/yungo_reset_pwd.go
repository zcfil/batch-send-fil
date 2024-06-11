package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"share-profit/result"
	"share-profit/utils"
)

type Restpwd struct{
	OldPwd string `json:"oldPwd" binding:"required"`
	NewPwd string	`json:"newPwd" binding:"required"`
}

//修改密码
func (y *Yungo) ReSetPwd(ctx *gin.Context){

	var params Restpwd
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		ctx.JSON(200, result.Fail(err))
		ctx.Abort()
		return
	}
	ctx.Request.FormValue("newpass")
	user, _ := utils.GetSubject(y.Client, ctx)
	//原密码校验
	oldpwd := utils.EncodePassword(user.Username, params.OldPwd)
	if user.Password != oldpwd {
		ctx.JSON(200, result.Fail(errors.New("密码错误")))
		ctx.Abort()
		return
	}
	//新密码替换
	newpwd := utils.EncodePassword(user.Username, params.NewPwd)
	//user.Password=newpwd
	//fmt.Println("id:",user.Id)
	//_, err = y.Engine.Id(user.Id).Cols("password").Update(user)
	//if err != nil {
	//	ctx.JSON(200, result.Fail(err))
	//	ctx.Abort()
	//	return
	//}
	err = y.UpdatePassowrd(newpwd,strconv.FormatInt(user.Id,10))
	if err != nil {
		ctx.JSON(200, result.Fail(err))
		ctx.Abort()
		return
	}
	//修改成功
	ctx.JSON(200,result.Ok(nil))
}

//个人信息
func (y *Yungo) MemberInfo(ctx *gin.Context){

	user, err := utils.GetSubject(y.Client, ctx)
	if err != nil {
		ctx.JSON(200, result.Fail(err))
		ctx.Abort()
		return
	}

	info := struct {
		Id         int64
		Username   string
		RealName   string
		HeadUrl    string
		Gender     int
		Email      string
		Mobile     string
		DeptId     int64
		SuperAdmin int
		Verified   int
		Status     int
		Remark     string
		DelFlag    int
		Creator    int64
		CreateDate time.Time
		Updater    int64
		UpdateDate time.Time
	}{
		Id:         user.Id,
		Username:   user.Username,
		RealName:   user.RealName,
		HeadUrl:    user.HeadUrl,
		Gender:     user.Gender,
		Email:      user.Email,
		Mobile:     user.Mobile,
		DeptId:     user.DeptId,
		SuperAdmin: user.SuperAdmin,
		Verified:   user.Verified,
		Status:     user.Status,
		Remark:     user.Remark,
		DelFlag:    user.DelFlag,
		Creator:    user.Creator,
		CreateDate: user.CreateDate,
		Updater:    user.Updater,
		UpdateDate: user.UpdateDate,
	}

	ctx.JSON(200,result.Ok(info))
}

