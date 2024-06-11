package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"share-profit/conf"
	"share-profit/entity"
)

//GetSubject 获取当前用户信息,直接从缓存中获取
func GetSubject(client *redis.Client, ctx *gin.Context) (*entity.SysUser, error) {
	var tokenStr string
	tokenStr = ctx.GetHeader(conf.API_KEY)
	if tokenStr ==""{
		tokenStr = ctx.DefaultQuery("token", "")
	}
	//从缓存中取到当前登录人信息
	userInfoBytes, err := client.Get(conf.TOKEN_PREFIX + tokenStr).Bytes()
	if err != nil {
		return nil, err
	}
	var sysUser = new(entity.SysUser)
	err = json.Unmarshal(userInfoBytes, sysUser)
	return sysUser, err
}

//直接用tokeen串从缓存拿
func GetSubjectByTokenStr(tokenStr string, client *redis.Client) (*entity.SysUser, error) {

	//从缓存中取到当前登录人信息
	userInfoBytes, err := client.Get(conf.TOKEN_PREFIX + tokenStr).Bytes()
	if err != nil {
		return nil, err
	}
	var sysUser = new(entity.SysUser)
	err = json.Unmarshal(userInfoBytes, sysUser)
	return sysUser, err

}

