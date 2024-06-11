package controller

import (
	"fmt"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"share-profit/entity"
	"share-profit/result"
	"share-profit/service"
	"share-profit/utils"
	"strings"
	"time"
)

func init() {

}

const (
	PledgeFIL32GiB = iota
	PledgeFIL1TiB
	PledgeFIL1PiB
)

type YungoLotus struct {
	*service.LotusService
}

func NewYungoLotus(s service.LotusService) *YungoLotus {

	return &YungoLotus{
		&s,
	}
}

const (
	PC2 = iota
	C2
)

func (l *YungoLotus) GetStatisticalFIL(ctx *gin.Context) {
	gas := l.StatisticalFIL()

	ctx.JSON(http.StatusOK, result.Ok(gas))
}

func (l *YungoLotus) GetPleageFIL(ctx *gin.Context) {
	gas := l.GetGas()
	ctx.JSON(http.StatusOK, result.Ok(gas))
}

//计算gas费用
func (l *YungoLotus) AddPc2AndC2(pc2, c2 *api.InvocResult, pledgeType int) *entity.YungoGas {
	gas := l.GetGasResul(pc2, c2)
	switch pledgeType {
	case PledgeFIL32GiB:
	case PledgeFIL1TiB:
		gas.Totalfil = gas.Totalfil * 32
		gas.Gasfil = gas.Gasfil * 32
		gas.Pleagefil = gas.Pleagefil * 32
	case PledgeFIL1PiB:
		gas.Totalfil = gas.Totalfil * 32 * 1024
		gas.Gasfil = gas.Gasfil * 32 * 1024
		gas.Pleagefil = gas.Pleagefil * 32 * 1024
	}
	return gas
}

//获取汇率
func (l *YungoLotus) GetCNYbyFIL(pc2, c2 *api.InvocResult, pledgeType int) *entity.YungoGas {
	gas := l.GetGasResul(pc2, c2)
	switch pledgeType {
	case PledgeFIL32GiB:
	case PledgeFIL1TiB:
		gas.Totalfil = gas.Totalfil * 32
		gas.Gasfil = gas.Gasfil * 32
		gas.Pleagefil = gas.Pleagefil * 32
	case PledgeFIL1PiB:
		gas.Totalfil = gas.Totalfil * 32 * 1024
		gas.Gasfil = gas.Gasfil * 32 * 1024
		gas.Pleagefil = gas.Pleagefil * 32 * 1024
	}
	return gas
}

func AddGasFIL(l *YungoLotus) {
	go func() {
		for {
			var ctx *gin.Context
			s, _ := l.Rpc.ChainHead(ctx)
			var index int64 = 5
			s1, _ := l.Rpc.ChainGetTipSetByHeight(ctx, s.Height()-abi.ChainEpoch(index), types.TipSetKey{})
			fmt.Println("height:", s.Height())
			fmt.Println("height:", s1.Height())

			ms, _ := l.Rpc.ChainGetBlockMessages(ctx, s1.Cids()[0])
			meskey := map[int][]*types.Message{}
			for _, val := range ms.BlsMessages {
				if val.Method == 6 {
					meskey[PC2] = append(meskey[PC2], val)
				}
				if val.Method == 7 {
					meskey[C2] = append(meskey[C2], val)
				}
			}
			var pc2, c2 *api.InvocResult

			for k, val := range meskey {
				for i := 0; i < len(val); i++ {
					rep, err := l.Rpc.StateReplay(ctx, types.TipSetKey{}, val[i].Cid())
					if err != nil {
						fmt.Println(err)
						continue
					}
					if rep == nil {
						fmt.Println("获取不到：", val[i].Cid())
						continue
					}
					if rep.GasCost.TotalCost.Int.Int64() == 0 {
						continue
					}
					if k == PC2 {
						if i == 0 {
							pc2 = rep
						}
						res, _ := new(big.Int).SetString(pc2.GasCost.TotalCost.Int.String(), 10)
						if !strings.Contains(res.Sub(res, rep.GasCost.TotalCost.Int).String(), "-") {
							pc2 = rep
							fmt.Println("cid:", rep.GasCost.Message, ",value:", rep.GasCost.TotalCost)
						}
					}
					if k == C2 {
						if i == 0 {
							c2 = rep
						}
						res, _ := new(big.Int).SetString(c2.GasCost.TotalCost.Int.String(), 10)
						if !strings.Contains(res.Sub(res, rep.GasCost.TotalCost.Int).String(), "-") {
							fmt.Println("cid:", c2.GasCost.Message, ",value:", rep.GasCost.TotalCost)
							c2 = rep
						}
					}
					//break
				}
			}
			//gas := l.AddPc2AndC2(pc2,c2,PledgeFIL32GiB)
			gas := l.GetGasResul(pc2, c2)
			fmt.Println(gas)
			val := *l.GetGas()
			if gas.Totalfil*2 < utils.DataToFloat64(val[0], "totalfil") || gas.Gasfil*2 < utils.DataToFloat64(val[0], "gasfil") || gas.Pleagefil*2 < utils.DataToFloat64(val[0], "pleagefil") {
				continue
			}
			if gas.Cnytofil == 0 {
				gas.Cnytofil = utils.DataToFloat64(val[0], "cnytofil")
			}
			l.Engine.Insert(gas)
			//ctx.JSON(http.StatusOK,result.Ok(gas))
			time.Sleep(time.Minute * 5)
		}
	}()
}
