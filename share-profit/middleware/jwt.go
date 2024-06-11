package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"log"
	"share-profit/conf"
	"share-profit/utils"
	"time"
)

var route = map[string]string{
	"/yungo/getVerificationCode":"",
	"/yungo/bindVerificationCode":"",
	"/yungo/verifyCode":"",
}
//JwtCheck 校验token是否生效的中间件
func JwtCheck(client *redis.Client) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		tokenStr := ctx.GetHeader(conf.API_KEY)
		//请求头是否携带token
		if tokenStr == "" {

			tokenquery := ctx.DefaultQuery("token", "")
			if tokenquery == ""{
				ctx.String(500, "%s不能为空", conf.API_KEY)
				ctx.Abort()
				return
			}
			tokenStr = tokenquery
		}
		fmt.Println("token:",tokenStr)
		//解析token
		token, err := jwt.ParseWithClaims(tokenStr, &utils.LoginClaims{}, func(token *jwt.Token) (i interface{}, e error) {

			return utils.SignKey, nil
		})
		//解析tokenStr报错,
		if err != nil && token == nil {
			log.Println("token解析报错:")
			ctx.String(500, "%s", err.Error())
			ctx.Abort()
			return
		}
		sys_user,err := utils.GetSubjectByTokenStr(tokenStr,client)
		if _,ok := route[ctx.Request.URL.Path];!ok{
			if sys_user.Verified!=1{
				log.Println("token未验证!:")
				ctx.String(500, "%s","token未验证！" )
				ctx.Abort()
				return
			}
		}

		//jwt字符串解析有效
		if _, ok := token.Claims.(*utils.LoginClaims); ok && token.Valid {
			log.Println("jwt valid pass...")
			ctx.Next()
		} else {
			//校验token是否过期，则重新生成token，并更新redis缓存
			claims := token.Claims.(*utils.LoginClaims)
			if isExpired := claims.VerifyExpiresAt(time.Now().Unix(), true); isExpired == false {
				//查询缓存中是否存在token，存在则重置token的有效期,不存在则说明登录过期
				if exist, err := client.Exists(conf.TOKEN_PREFIX + tokenStr).Result(); exist == 0 || err != nil {
					//redis中不存在token
					ctx.String(500, "token失效,用户登录过期!")
					ctx.Abort()
					return
				} else {
					//缓存中还存在token,将缓token有效时间重置
					client.Expire(conf.TOKEN_PREFIX+tokenStr, 2*time.Second*time.Duration(conf.TOKEN_EFFECT_TIME))

					ctx.Next()
				}
			}
		
		}

	}

}


