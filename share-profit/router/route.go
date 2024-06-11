package router

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/go-xorm/xorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"share-profit/conf"
	"share-profit/controller"
	"share-profit/middleware"
	"share-profit/result"
	"syscall"
)

func Serve(client *redis.Client,
	engine *xorm.Engine,
	enforcer *casbin.Enforcer,
	repo *conf.Repo,
	pro *conf.Project,
	user controller.User,
	login *controller.Login,
	dept controller.Dept,
	admin controller.SuperAdmin,
	wallet *controller.Wallet,
	yungo controller.Yungo,
	lotus *controller.YungoLotus,
	fi *controller.Finance) {

	r := gin.Default()
	//cors跨域配置
	r.Use(conf.Cors())
	//文件大小限制
	r.MaxMultipartMemory = repo.MaxSize
	//云构
	yungoGroup := r.Group("/yungo")

	yungoGroup.Use(middleware.JwtCheck(client))
	{
		yungoGroup.POST("/upload", yungo.Upload)
		yungoGroup.POST("/wallet/default", yungo.GetWallet)
		yungoGroup.GET("/list/imports", yungo.ListImportFiles)
		yungoGroup.GET("/queryask", yungo.ClientQueryAsk)
		yungoGroup.POST("/deal", yungo.ClientDeal)
		yungoGroup.GET("/dealinfo", yungo.DealInfo)
		yungoGroup.GET("/ftypes", yungo.ListFileTypes)
		//重置密码
		yungoGroup.POST("/resetpwd", yungo.ReSetPwd)
		// 个人信息
		yungoGroup.GET("/memberinfo", yungo.MemberInfo)
		yungoGroup.GET("/download", yungo.DownLoad)

		yungoGroup.POST("/add", func(ctx *gin.Context) {
			ctx.JSON(200, result.Ok(""))
		})
		//钱包接口
		//yungoGroup.POST("/send",wallet.WalletSend) //转账
		yungoGroup.GET("/getMnemonic", wallet.GetMnemonic) //获取助记词
		yungoGroup.POST("/newWallet", wallet.NewWallet)    //创建钱包地址
		yungoGroup.GET("/walletList", wallet.WalletList)   //获取钱包列表
		yungoGroup.POST("/setWallet", wallet.SetWallet)    //设置为默认钱包
		yungoGroup.POST("/delWallet", wallet.DelWallet)    //删除钱包
		//yungoGroup.POST("/importWallet", wallet.ImportWallet)     //通过私钥导入钱包
		//yungoGroup.GET("/exportPrivateKey", wallet.ExportWallet)  //导出私钥
		//yungoGroup.POST("/importMnemonic", wallet.ImportMnemonic) //通过助记词导入钱包
		//yungoGroup.GET("/exportMnemonic", wallet.ExportMnemonic)  //导出助记词

		//2021年3月9日10:43:55  分币系统需求
		yungoGroup.POST("/uploadApply", fi.Upload)            //导入提币申请
		yungoGroup.GET("/getApplyList", fi.ApplyList)         //获取提币列表
		yungoGroup.POST("/sends", fi.WalletSends)             //批量转账
		yungoGroup.POST("/refuse", fi.UpdateStatus)           //批量拒绝
		yungoGroup.GET("/walletBalance", fi.WalletAndBalance) //钱包余额
		yungoGroup.GET("/getBatchList", fi.BatchList)         //获取批次列表
		yungoGroup.POST("/batchSends", fi.BatchSends)         //按批次转账
		yungoGroup.POST("/batchRefuse", fi.BatchRefuse)       //按批次拒绝
		yungoGroup.GET("/getConfig", fi.GetConfig)            //获取配置列表
		yungoGroup.POST("/setConfig", fi.SetConfig)           //设置配置

		yungoGroup.GET("/getManualList", fi.ManualList)   //获取手动输入列表
		yungoGroup.POST("/updateManual", fi.UpdateManual) //批量手动列表拒绝
		yungoGroup.POST("/sendManual", fi.SendManual)     //批量手动列表转账
		yungoGroup.POST("/manualAdd", fi.ManualAdd)       //添加

		yungoGroup.GET("/getVerificationCode", login.CreateSecret) //创建谷歌验证码
		yungoGroup.POST("/bindVerificationCode", login.BindCode)   //绑定谷歌验证
		yungoGroup.POST("/verifyCode", login.VerifyCode)           //验证谷歌验证

	}
	//无权限接口  getVerificationCode  bindVerificationCode  verifyCode
	//登录
	r.POST("/login", login.DoLogin)
	//注册
	r.POST("/register", login.DoRegister)
	//商户注册
	r.POST("/merchants/register", login.MerchantsRegister)

	r.GET("/pleagefil", lotus.GetPleageFIL)

	r.GET("/statisticalfil", lotus.GetStatisticalFIL)

	userGroup := r.Group("/user")
	userGroup.Use(middleware.JwtCheck(client))
	{
		userGroup.POST("/logout", user.DoLogout)
	}

	//组织（需要token校验，角色校验）
	deptGroup := r.Group("/dept")
	deptGroup.Use(middleware.JwtCheck(client), middleware.CheckPermission(client, enforcer))
	{
		deptGroup.POST("/add", dept.InsertDept)
		deptGroup.GET("/verify", dept.VerifyList)
		deptGroup.PUT("/check/:id", dept.CheckDept)
	}
	//系统功能
	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.JwtCheck(client), middleware.CheckAdmin(client))
	{
		adminGroup.POST("/addRole", admin.AddRole)
		adminGroup.POST("/register", login.DoRegister)
	}
	//swagger文档
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//log.Printf("swagger : %s", "http://localhost:3000/swagger/index.html")
	srv := &http.Server{
		Addr:    ":" + pro.Host,
		Handler: r,
	}
	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	//go func() {
	//	<-sign
	//	log.Println("Shutdown............ ...")
	//	if err := srv.Shutdown(context.Background()); err != nil {
	//		log.Fatal("Server Shutdown:", err)
	//	}
	//}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
	log.Println("Release connection resources")
	client.Close()
	engine.Close()
	yungo.YungoRpc.Closer()
	log.Println("Server exiting")

}
