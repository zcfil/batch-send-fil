package models

import (
	"github.com/go-redis/redis/v7"
	"github.com/go-xorm/xorm"
	"go.uber.org/dig"
	"share-profit/rpc"
	"share-profit/utils"
)

type LoginModels struct {
	dig.In
	*xorm.Engine
	*redis.Client
	*rpc.YungoRpc
}

func (l *LoginModels)GetUserPwd(phone string)string{
	sql := `select password from user where phone='`+phone+ `'`
	res,_ := l.Engine.QueryString(sql)
	return res[0]["password"]
}
func (l *LoginModels)GetUserCountByPhone(phone string)string{
	sql := `select count(1) total from user where phone='`+phone+ `'`
	res,_ := l.Engine.QueryString(sql)
	return res[0]["total"]
}

func (l *LoginModels)UpdateSysUser(param map[string]string)error{
	sql := `update sys_user set verification =:verification where id =:id `
	sql = utils.SqlReplaceParames(sql,param)
	_,err := l.Exec(sql)

	return err
}

func (l *LoginModels)GetSysUserVerification(id string)string{
	sql := `select verification from sys_user where id='`+id+ `'`
	res,err := l.Engine.QueryString(sql)
	if err !=nil||len(res)==0{
		return ""
	}

	return res[0]["verification"]
}