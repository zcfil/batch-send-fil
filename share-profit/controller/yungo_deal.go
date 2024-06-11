package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"
	"log"
	"share-profit/entity"
	"share-profit/result"
)

//交易参数
type Deal struct {
	Id int
	MinerId string
	Price string
}

//查询所有miner
func (y *Yungo) ClientQueryAsk(ctx *gin.Context) {

	minerStr := ctx.DefaultQuery("miner","")
	if minerStr == "" {
		log.Println("lotus client local:", errors.New("minerID 是空字符串"))
		ctx.JSON(200, result.Fail(errors.New("minerID 是空字符串")))
		ctx.Abort()
		return
	}
	//miner
	miner, err := address.NewFromString(minerStr)
	if err != nil {
		log.Println("lotus client local:", err)
		ctx.JSON(200, result.Fail(err))
		ctx.Abort()
		return
	}
	//询价peer.ID
	log.Println("miner:", miner)
	mi, err := y.Rpc.StateMinerInfo(context.TODO(), miner, types.EmptyTSK)
	if err != nil {
		log.Println("lotus client local:", err)
		ctx.JSON(200, result.Fail(err))
		ctx.Abort()
		return
	}
	//lotus client query-ask
	storageAsk, err := y.Rpc.ClientQueryAsk(context.TODO(), *mi.PeerId, miner)
	if err != nil {
		log.Println("lotus client local:", err)
		ctx.JSON(200, result.Fail(err))
		ctx.Abort()
		return
	}
	ctx.JSON(200, result.Ok(storageAsk))
}

//发起交易
func (y *Yungo) ClientDeal(ctx *gin.Context) {

	var deal Deal
	err := ctx.ShouldBindJSON(&deal)
	if err !=nil {
		log.Println("参数错误:", err)
		ctx.JSON(200, result.Fail(err))
		ctx.Abort()
		return
	}

	file := new(entity.YungoFile)
	has, err := y.Engine.Where("id=?", deal.Id).Get(file)
	if err != nil  {
		log.Println("lotus client deal:", err)
		ctx.JSON(200, result.Fail(err))
		ctx.Abort()
		return
	}
	if  has == false {
		ctx.JSON(200, result.Fail(errors.New("数据记录不存在")))
		ctx.Abort()
		return
	}
	//miner
	miner, err := address.NewFromString(deal.MinerId)
	//默认钱包地址
	wallet, _ := y.Rpc.WalletDefaultAddress(context.TODO())
	//root cid
	rootCid, err := cid.Decode(file.RootCid)
	if err != nil {
		log.Println("lotus client deal:", err)
		ctx.JSON(200, result.Fail(err))
		ctx.Abort()
		return
	}
	dealCid, err := y.Rpc.ClientStartDeal(context.TODO(), &api.StartDealParams{
		Data: &storagemarket.DataRef{
			TransferType: storagemarket.TTGraphsync,
			Root:         rootCid,
		},
		Wallet:             wallet,
		Miner:              miner,
		EpochPrice:         types.NewInt(100000000000),
		MinBlocksDuration:  518400,
		ProviderCollateral: big.Int{},
		//ProviderCollateral: big.Int{},
		//DealStartEpoch: 0,
		//FastRetrieval:  false,
		VerifiedDeal: false,
	})
	if err != nil {
		log.Println("lotus client deal:", err)
	}
	file.DealId= dealCid.String()
	//affected, err := engine.Id(id).Update(user)
	affected, err := y.Engine.ID(file.Id).Update(file)
	if err != nil {
		log.Println("update dealID:", err)
		ctx.JSON(200, result.Fail(err))
		ctx.Abort()
		return
	}
	fmt.Println("更新结果:",affected)
	ctx.JSON(200, result.Ok(dealCid))
}

//查询交易状态
func (y *Yungo) DealInfo(ctx *gin.Context) {
	dealID  := ctx.Query("deal_id")
	if dealID == ""{
		log.Println("dealID 不能为空:")
		ctx.JSON(200, result.Fail(errors.New("dealID不能为空")))
		ctx.Abort()
		return
	}
	proserCid, err := cid.Decode(dealID)
	if err != nil {
		log.Println("生成CID:", err)
		ctx.JSON(200, result.Fail(err))
		ctx.Abort()
		return
	}

	info, err := y.Rpc.ClientGetDealInfo(context.TODO(), proserCid)
	if err != nil {
		log.Println("deal info:", err)
		ctx.JSON(200, result.Fail(err))
		ctx.Abort()
		return
	}

	ctx.JSON(200,result.Ok(info))
}
