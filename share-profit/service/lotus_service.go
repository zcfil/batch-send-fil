package service

import (
	"fmt"
	"github.com/filecoin-project/lotus/api"
	"github.com/go-redis/redis/v7"
	"github.com/go-xorm/xorm"
	"go.uber.org/dig"
	"math/big"
	"regexp"
	"share-profit/entity"
	"share-profit/rpc"
	"share-profit/utils"
	"strconv"
	"time"
)

type LotusService struct {
	dig.In
	*redis.Client
	*rpc.YungoRpc
	*xorm.Engine
}

func (l *LotusService) GetGasFIL(pc2, c2 *api.InvocResult) (result float64) {
	//pc2gas费
	f099, _ := new(big.Int).SetString(pc2.GasCost.TotalCost.Int.String(), 10)
	total := f099.Add(pc2.GasCost.TotalCost.Int, pc2.GasCost.MinerTip.Int)
	fmt.Println("f099,pc2,total:", f099, pc2.GasCost.TotalCost.Int, total)

	//C2gas费
	f099, _ = new(big.Int).SetString(c2.GasCost.TotalCost.Int.String(), 10)
	totalc2 := f099.Add(c2.GasCost.TotalCost.Int, c2.GasCost.MinerTip.Int)
	fmt.Println("f099,pc2,totalc2:", f099, c2.GasCost.TotalCost.Int, totalc2)

	count := total.Add(total, totalc2).String()
	fmt.Println("total,totalc2,count:", total, totalc2, count)
	result = utils.NanoOrAttoToFIL(count, utils.AttoFIL)

	return
}

//质押费
func (l *LotusService) GetPleageFIL(pc2, c2 *api.InvocResult) (result float64) {

	pc2pledge, _ := new(big.Int).SetString(pc2.Msg.Value.Int.String(), 10)
	count := pc2pledge.Add(pc2pledge, c2.Msg.Value.Int).String()
	result = utils.NanoOrAttoToFIL(count, utils.AttoFIL)
	return
}

func (l *LotusService) GetTotalFIL(pc2, c2 *api.InvocResult) (result float64) {
	//pc2gas费
	f099, _ := new(big.Int).SetString(pc2.GasCost.TotalCost.Int.String(), 10)
	total := f099.Add(f099, pc2.GasCost.MinerTip.Int)
	//C2gas费
	f099, _ = new(big.Int).SetString(c2.GasCost.TotalCost.Int.String(), 10)
	totalc2 := f099.Add(f099, c2.GasCost.MinerTip.Int)

	gascount := total.Add(total, totalc2)

	//质押费
	pc2pledge, _ := new(big.Int).SetString(pc2.Msg.Value.Int.String(), 10)
	pledgecount := pc2pledge.Add(pc2pledge, c2.Msg.Value.Int)

	//总和
	count := gascount.Add(gascount, pledgecount).String()
	result = utils.NanoOrAttoToFIL(count, utils.AttoFIL)

	return
}

func (l *LotusService) GetGasResul(pc2, c2 *api.InvocResult) *entity.YungoGas {

	return &entity.YungoGas{
		l.GetGasFIL(pc2, c2),
		l.GetPleageFIL(pc2, c2),
		l.GetTotalFIL(pc2, c2),
		GetCNYbyFIL(),
		time.Now().Unix(),
	}
}
func (l *LotusService) GetGas() *[]map[string]string {
	sql := ` select * from yungo_gas group by gas_id desc limit 0,5 `
	//var mp []map[string][]byte
	mp, err := l.Engine.QueryString(sql)
	if err != nil {
		fmt.Errorf("%v", err)
	}
	if mp == nil {
		fmt.Println("GetGas无数据")
		return nil
	}
	fmt.Println("云构：", mp[0]["pleagefil"])
	//gas.Cnytofil,_ =  strconv.ParseFloat(string(mp[0]["cnytofil"]),64)
	//gas.Gasfil,_ =  strconv.ParseFloat(string(mp[0]["gasfil"]),64)
	//gas.Totalfil,_ = strconv.ParseFloat(string(mp[0]["totalfil"]),64)

	return &mp
}

//获取统计数据
func (l *LotusService) StatisticalFIL() *[]map[string]string {
	sql := `select * from (select round(avg(gasfil),4) gasfil,round(avg(gas_id),0) gas_id,round(avg(pleagefil),4) pleagefil,round(avg(totalfil),4) totalfil,round(avg(cnytofil),4) cnytofil,
            FROM_UNIXTIME(create_time,'%m/%d') create_time
            from yungo_gas
			group by FROM_UNIXTIME(create_time,'%m/%d')
			order by create_time desc limit 0,7 )a order by create_time  `
	//var mp []map[string][]byte
	mp, err := l.Engine.QueryString(sql)
	if err != nil {
		fmt.Errorf("%v", err)
	}

	return &mp
}

const (
	//CNYReg = `<span data-v-34fdfde4>≈¥([0-9]*\.[0-9]+)`
	CNYReg = `"legal_currency_price":([0-9]*\.[0-9]+)`
)

func (l *LotusService) InsertGas(gas entity.YungoGas) {
	l.Engine.Insert(&gas)
	return
}

const URL = "https://www.mytokencap.com/currency/fil/821765876"

func GetCNYbyFIL() float64 {
	str := utils.Get(URL)
	mc := regexp.MustCompile(CNYReg)
	submatch := mc.FindAllStringSubmatch(str, -1)
	var fil float64
	for _, m := range submatch {
		fil, _ = strconv.ParseFloat(m[1], 64)
		fmt.Println("云构：", m[1], fil)
		break
	}
	return fil
}
