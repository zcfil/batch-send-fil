package main

import (
	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
	"go.uber.org/dig"
	"log"
	"os"
	"share-profit/conf"
	"share-profit/controller"
	"share-profit/db"
	_ "share-profit/docs"
	"share-profit/router"
	"share-profit/rpc"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email 1503780117@qq.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 183.61.251.226:3000
// @BasePath
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name token
func main() {

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Value: "./conf/config.toml",
				Usage: "启动的配置文件路径格式为toml",
			},
			&cli.StringFlag{
				Name:  "geo",
				Value: "./conf/GeoLite2-City.mmdb",
				Usage: "GeoIP数据库文件路径,将IP地址解析出地理位置（国家，城市，经纬度）",
			},
		},
		Action: func(c *cli.Context) error {

			var config conf.Config
			if _, err := toml.DecodeFile(c.String("path"), &config); err != nil {
				// handle error
				return err
			}
			container := dig.New()
			container.Provide(func() *conf.Repo {
				return &conf.Repo{
					config.Repo.UploadPath,
					config.Repo.MaxSize,
				}
			})
			container.Provide(func() *conf.Project {
				return &conf.Project{
					config.Project.Host,
				}
			})
			container.Provide(db.NewSqlEngine(config.Mysql))
			container.Provide(conf.NewRedisClient(config.Redis))
			container.Provide(rpc.NewYungoRpc(config.Lotus))
			container.Provide(controller.NewLogin)
			container.Provide(controller.NewYungoLotus)
			container.Provide(controller.NewFinance)
			//geoDB配置路径
			container.Provide(controller.NewWallet)
			container.Provide(conf.LoadCasbin(config.Mysql))
			//err = container.Invoke(controller.AddGasFIL) //会崩溃
			//err = container.Invoke(controller.GetGasCost)
			err := container.Invoke(router.Serve)
			//go func(){
			//	var ctx *gin.Context
			//	//var l controller.YungoLotus
			//	service.AddGasFIL(ctx)
			//	time.Sleep(time.Minute * 5)
			//}()
			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("启动失败:", err)
	}

}
