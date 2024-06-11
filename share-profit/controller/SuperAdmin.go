package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/go-xorm/xorm"
	"go.uber.org/dig"
	"net/http"
	"share-profit/entity"
	"share-profit/result"
	"share-profit/utils"
)

type SuperAdmin struct {
	dig.In
	*xorm.Engine
	*redis.Client
}

//@description 新增角色
//@accept json
//@Param SysRole body entity.SysRole true "角色"
//@Success 200 {object} gin.H
//@router /admin/addRole [post]
//@Security ApiKeyAuth
func (admin *SuperAdmin) AddRole(ctx *gin.Context) {

	var sysRole = new(entity.SysRole)

	err := ctx.ShouldBindJSON(sysRole)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		ctx.Abort()
		return
	}
	//生成id
	sysRole.Id = utils.Node().Generate().Int64()
	user, _ := utils.GetSubject(admin.Client,ctx)
	sysRole.Creator = user.Id

	if i, err := admin.Insert(sysRole); err == nil && i == 1 {

		ctx.JSON(200, result.Ok(sysRole))
	} else {

		ctx.JSON(200, result.Fail(err))
	}

}
