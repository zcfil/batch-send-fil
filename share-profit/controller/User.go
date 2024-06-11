package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/go-xorm/xorm"
	"go.uber.org/dig"
	"log"
	"share-profit/conf"
)


type User struct {
	dig.In
	*xorm.Engine
	*redis.Client
}

//@description 用户退出
//@accept json
//@Success 200 {object} gin.H
//@router /user/logout [post]
//@Security ApiKeyAuth
func (user *User) DoLogout(ctx *gin.Context) {

	tokenStr := ctx.GetHeader(conf.API_KEY)
	//检查是否携带token
	if tokenStr == "" {

		ctx.String(500, "token为空!")
		ctx.Abort()
		return
	}
	//查询token是否存在缓存里
	if exist, err := user.Exists(conf.TOKEN_PREFIX + tokenStr).Result(); exist == 0 || err != nil {

		log.Println("exist:", exist)
		ctx.String(500, "token已过期!")
		ctx.Abort()
		return
	}

	if result, err := user.Del(conf.TOKEN_PREFIX + tokenStr).Result(); result == 0 || err != nil {

		log.Println("result:", result)
		ctx.String(500, "退出失败:%v",err)
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"code": 0, "data": nil, "msg": "success"})
}


