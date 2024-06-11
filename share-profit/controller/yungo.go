package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/filecoin-project/lotus/api"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/go-xorm/xorm"
	"go.uber.org/dig"
	"log"
	"net/http"
	"strconv"
	"time"
	"share-profit/conf"
	"share-profit/entity"
	"share-profit/result"
	"share-profit/rpc"
	"share-profit/utils"
)

type Yungo struct {
	dig.In
	*rpc.YungoRpc
	*xorm.Engine
	*redis.Client
	*conf.Repo
}


//文件上传
func (y *Yungo) Upload(ctx *gin.Context) {
	// single file
	file, err := ctx.FormFile("file")
	ftype := ctx.DefaultPostForm("ftype", "5")
	ft, err := strconv.ParseInt(ftype, 10, 32)
	if err != nil {
		ctx.JSON(200, result.Fail(err))
		ctx.Abort()
		return
	}
	fmt.Println("ftype----------",ftype)
	if err != nil {
		log.Println(err)
		ctx.String(http.StatusOK, err.Error())
		return
	}
	log.Println("file-name:", file.Filename)
	// Upload the file to specific dst.
	ctx.SaveUploadedFile(file, y.UploadPath+file.Filename)
	//lotus client import
	res, err := y.Rpc.ClientImport(context.TODO(), api.FileRef{
		Path:  y.Repo.UploadPath + file.Filename,
		IsCAR: false,
	})
	if err != nil {
		ctx.JSON(200, err.Error())
		ctx.Abort()
		return
	}

	user, err := utils.GetSubject(y.Client,ctx)
	if err != nil {
		ctx.JSON(200, err.Error())
		ctx.Abort()
		return
	}
	//存到databae
	yfile := &entity.YungoFile{
		MemberId:   user.Id,
		Type:       int(ft),
		Filename:   file.Filename,
		FileSize:   strconv.FormatInt(file.Size, 10),
		Price:      "0.00000001",
		RootCid:    res.Root.String(),
		MinerId:    "",
		DealId:     "",
		CreateTime: time.Now().UTC(),
		UpdateTime: time.Now().UTC(),
	}
	_, err = y.Insert(yfile)

	if err != nil {

		ctx.JSON(http.StatusBadRequest, err.Error())
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, result.Ok(yfile))
}

func (y *Yungo) GetWallet(ctx *gin.Context) {

	address, err := y.Rpc.WalletDefaultAddress(context.Background())

	if err != nil {
		log.Fatal("lotus wallet default:", err)
		ctx.String(500, err.Error())
		return
	}
	res, err := y.Rpc.ClientImport(context.TODO(), api.FileRef{
		Path:  y.UploadPath + "casbin-express.rar",
		IsCAR: false,
	})
	if err != nil {
		log.Println("lotus client import:", err)
		ctx.String(500, err.Error())
		return
	}
	fmt.Println("RootCid:", res.Root)
	ctx.JSON(200, address.String())
}

//分页查询import的文件
func (y *Yungo) ListImportFiles(ctx *gin.Context) {

	pageNum := ctx.DefaultQuery("pageNum", "1")
	pageSize := ctx.DefaultQuery("pageSize", "10")
	//文件类型
	fileType := ctx.DefaultQuery("type", "")
	pageNumber, err := strconv.Atoi(pageNum)
	if err != nil {
		ctx.JSON(200,result.Fail(errors.New("pageNum参数有误")))
		ctx.Abort()
		return
	}
	pageSizeNumber, err := strconv.Atoi(pageSize)
	if err != nil {
		ctx.JSON(200,result.Fail(errors.New("pageSize参数有误")))
		ctx.Abort()
		return
	}
	//获取用户info
	user, err := utils.GetSubject(y.Client,ctx)
	if err != nil {
		ctx.JSON(200, err.Error())
		ctx.Abort()
		return
	}
	fmt.Println("user-info:", user)
	importFiles := make([]entity.YungoFile, 0)
	var sql string
	if fileType =="" {
		sql =fmt.Sprintf("member_id = %d",user.Id)
	}else {
		sql = fmt.Sprintf("member_id = %d AND type =%s ",user.Id,fileType)
	}
	err = y.Engine.Where(sql).Limit(pageSizeNumber,(pageNumber-1)*pageSizeNumber).Find(&importFiles)
	if err != nil {
		log.Println("lotus client local:", err)
		ctx.JSON(200, result.Fail(err))
		ctx.Abort()
		return
	}
	importF := new(entity.YungoFile)
	total, _ := y.Engine.Where(sql).Count(importF)
	res := utils.NewPagination(int(pageNumber),int(pageSizeNumber),int(total),importFiles)
	ctx.JSON(200, result.Ok(res))
}

//文件类型列表
func (y *Yungo) ListFileTypes(ctx *gin.Context) {

	var types  []entity.YungoFileType
	err := y.Find(&types)
	if err !=nil {
		log.Println("lotus client local:", err)
		ctx.JSON(200, result.Fail(err))
		ctx.Abort()
		return
	}
	ctx.JSON(200,result.Ok(types))
}

func (y *Yungo) DownLoad(ctx *gin.Context) {

	fid := ctx.DefaultQuery("id","")
	if fid ==""{
		ctx.JSON(200, result.Fail(errors.New("参数不能为空")))
		ctx.Abort()
		return
	}

	user, err := utils.GetSubject(y.Client, ctx)
	if err != nil {
		ctx.JSON(200, result.Fail(err))
		ctx.Abort()
		return
	}
	var file entity.YungoFile
	has, err := y.Engine.Where("id = ? AND member_id = ?", fid,user.Id).Get(&file)
	if err != nil {
		ctx.JSON(200, result.Fail(err))
		ctx.Abort()
		return
	}
	if has {
		//查询文件
		ctx.File(fmt.Sprintf("%s%s",y.Repo.UploadPath,file.Filename))
	}else {
		ctx.JSON(200, result.Fail(errors.New("文件不存在")))
		ctx.Abort()
		return
	}
}

func (y *Yungo)UpdatePassowrd(pwd string,id string)error{
	sql := `update sys_user set password='` +pwd+"',update_date=now() where id = "+id
	_,err := y.Engine.Exec(sql)
	return err
}