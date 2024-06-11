package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/go-redis/redis/v7"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	"go.uber.org/dig"
	"log"
	"math/rand"
	"mime/multipart"
	"share-profit/conf"
	"share-profit/entity"
	"share-profit/rpc"
	"share-profit/utils"
	"strconv"
	"strings"
	"time"
)

type FinanceService struct {
	dig.In
	*redis.Client
	*rpc.YungoRpc
	*xorm.Engine
}

func CreateCaptcha() string {
	return fmt.Sprintf("%08v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000000))
}
func (f *FinanceService) GetDics(str string) map[string]string {
	q1 := ` select dict_name,dict_value from sys_dict where dict_type ='` + str + `' `
	dics, _ := f.Engine.QueryString(q1)
	dickey := make(map[string]string)
	for _, val := range dics {
		key := ""
		value := ""
		for k, v := range val {
			if k == "dict_name" {
				key = v
			}
			if k == "dict_value" {
				value = v
			}
		}
		dickey[key] = value
	}
	return dickey
}

//导入
func (f *FinanceService) UploadApply(file multipart.File, Size int64, user *entity.SysUser) error {
	buf := make([]byte, Size)
	n, _ := file.Read(buf)

	xf, _ := xlsx.OpenBinary(buf[:n])
	//获取字典
	dickey := f.GetDics("audit")
	statuskey := f.GetDics("status")
	//获取配置
	param := make(map[string]string)
	admin := strconv.FormatInt(user.AdminId, 10)
	param["username"] = user.Username
	param["admin_id"] = admin
	param["c_name"] = "charge"

	//charge,_ := strconv.ParseFloat(config,64)

	str := ",(:admin_id,:bnum"
	sql := `insert into withdraw(admin_id,batch_count,`
	sql1 := `select count(1)+1 bnum from batchnum a where admin_id=:admin_id`

	sql1 = utils.SqlReplaceParames(sql1, param)

	batch, err := f.Engine.QueryString(sql1)
	if err != nil {
		return err
	}

	param["bnum"] = batch[0]["bnum"]

	for _, sheet := range xf.Sheets {
		if len(sheet.Rows) == 0 {
			continue
		}
		//cls := len(sheet.Rows[0].Cells)
		exist := make([]string, len(sheet.Rows[0].Cells))
		cs := 0
		for j, row := range sheet.Rows {
			//value := row.Cells[1].String()+row.Cells[0].String()
			//if value == "" || value == "合计"{
			//	sql = sql[0:len(sql)-len(str)]
			//	break
			//}
			value := ""
			for i := 0; i < len(row.Cells); i++ {
				value += row.Cells[i].String()
			}
			if strings.Contains(value, "合计") || len(value) < len(row.Cells) {
				if j > 1 {
					sql = sql[:len(sql)-len(str)]
				}
				break
			}

			for i, cell := range row.Cells {
				if j == 0 {
					if _, ok := dickey[cell.String()]; ok {
						exist[i] = cell.String()
						sql += dickey[cell.String()] + ","
					}
					if i == len(row.Cells)-1 {
						sql = sql[0 : len(sql)-1]
						//sql += ")values(:admin_id,(select count(1)+1 a from batchnum a where admin_id=:admin_id),"
						sql += ")values(:admin_id,:bnum"
						cs = len(row.Cells) - 1
					}
				} else if len(exist) > i {
					if len(exist[i]) > 0 {
						value := cell.String()
						if exist[i] == "状态" {
							value = statuskey[value]
							if value == "" {
								value = "0"
							}
						}
						sql += ",'" + value + "'"
					}
					if i == len(row.Cells)-1 || i == cs {
						//sql = sql[:len(sql)-1]
						sql += ")"
					}
				}
			}

			if j >= 1 && j != len(sheet.Rows)-1 {
				sql += str
			}
		}
	}
	//if sql[len(sql)-1] == ',' {
	//	sql = sql[:len(sql)-1] + ")"
	//}
	sess := f.Engine.NewSession()
	sess.Begin()

	sql = utils.SqlReplaceParames(sql, param)

	re, err := sess.Exec(sql)
	defer func() {
		if err != nil {
			sess.Rollback()
			return
		}
		sess.Commit()
	}()
	if err != nil {
		fmt.Println("导入错误2：", err)
		return err
	}
	count, _ := re.RowsAffected()
	//bsql := ` insert into batchnum(operator, num, admin_id,count,create_time)value(:username,(select count(1)+1 a from batchnum a where admin_id=:admin_id) ,:admin_id,:count,now()); `
	bsql := ` insert into batchnum(operator, num, admin_id,count,create_time)value(:username,:bnum ,:admin_id,:count,now()); `

	param["count"] = strconv.FormatInt(count, 10)
	bsql = utils.SqlReplaceParames(bsql, param)

	_, err = sess.Exec(bsql)
	if err != nil {
		return err
	}

	return nil
}
func (f *FinanceService) GetApplyList(param map[string]string) ([]map[string]string, error) {
	str := ""
	if param["status"] != "" {
		str += " and status=" + param["status"]
	}
	if param["num"] == "" {
		param["num"] = `''`
	}
	sql := ` select *,if(w.filcount=0,0,w.amount-w.filcount) charge1 from withdraw w where admin_id=:admin_id and
            batch_count in (if(` + param["num"] + `='',
				(select max(batch_count) from withdraw a 
				where admin_id= :admin_id ` + str + ` ),` + param["num"] + `))
			` + str
	param["sort"] = "create_time"
	param["order"] = "desc"
	sql = utils.SqlReplaceParames(sql, param)
	//求总数
	var err error
	param["total"], err = utils.GetTotalCount(f.Engine, sql)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	sql += utils.LimitAndOrderBy(param)
	mp, err := f.Engine.QueryString(sql)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return mp, err
}

//手动输入列表查询
func (f *FinanceService) GetManualList(param map[string]string) ([]map[string]string, error) {
	str := ""
	if param["status"] != "" {
		str += " and status=" + param["status"]
	}
	sql := ` select * from withdraw_1 where admin_id=:admin_id ` + str
	param["sort"] = "create_time"
	param["order"] = "desc"
	sql = utils.SqlReplaceParames(sql, param)
	//求总数
	var err error
	param["total"], err = utils.GetTotalCount(f.Engine, sql)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	sql += utils.LimitAndOrderBy(param)
	mp, err := f.Engine.QueryString(sql)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return mp, err
}

//修改手动列表转账信息
func (f *FinanceService) UpdateManualInfo(param map[string]string) error {

	sql := ` update withdraw_1 set cid=:cid,filcount =:value,address=:to, status = :status,update_time=now(),errormsg=:error where id=:id `

	sql = utils.SqlReplaceParames(sql, param)
	_, err := f.Engine.Exec(sql)
	if err != nil {
		log.Println("withdraw_1状态未改变：", param["id"])
		return err
	}
	return nil
}

//修改手动列表状态
func (f *FinanceService) UpdateManualList(ids, status string) error {
	if ids == "" || status == "" {
		return errors.New("参数不能为空！")
	}
	id := strings.ReplaceAll(ids, ",", "','")
	sql := ` update withdraw_1 set status = ` + status + `,update_time=now() where id in('` + id + `') `
	_, err := f.Engine.Exec(sql)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (f *FinanceService) ManualAdds(param map[string]string) error {
	param["user_id"] = CreateCaptcha()
	sql := ` insert into withdraw_1(user_id,user_name,amount,status,address,admin_id,coin_type,create_time )value(:user_id,:user_name,:amount,"0",:address,:admin_id,"FIL",now())`
	sql = utils.SqlReplaceParames(sql, param)
	_, err := f.Engine.Exec(sql)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (f *FinanceService) Sends(to string, fil float64, user *entity.SysUser) (*types.SignedMessage, error) {

	//钱包地址from
	from, err := address.NewFromString(f.DefaultAddress(strconv.FormatInt(user.AdminId, 10)))
	if err != nil {
		return nil, err
	}
	//钱包地址to
	addrTo, err := address.NewFromString(to)
	if err != nil {
		fmt.Println("1 ", err, to)
		return nil, err
	}
	val, err := types.ParseFIL(strconv.FormatFloat(fil, 'f', -1, 64))
	if err != nil {
		return nil, fmt.Errorf("failed to parse amount: %w", err)
	}
	sm, err := f.Rpc.MpoolPushMessage(context.TODO(), &types.Message{
		To:    addrTo,
		From:  from,
		Value: types.BigInt(val),
		//GasLimit:   472835,
		//GasFeeCap:  abi.NewTokenAmount(101544),
		//GasPremium: abi.NewTokenAmount(100490),
	}, nil)
	if err != nil {
		return nil, err
	}
	//MSG <- sm.Cid()

	return sm, err
}

//修改状态
func (f *FinanceService) UpdateApplyList(ids, status string) error {
	if ids == "" || status == "" {
		return errors.New("参数不能为空！")
	}
	id := strings.ReplaceAll(ids, ",", "','")
	sql := ` update withdraw set status = ` + status + `,update_time=now() where id in('` + id + `') `
	_, err := f.Engine.Exec(sql)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//修改转账信息
func (f *FinanceService) UpdateTransferInfo(param map[string]string) error {

	sql := ` update withdraw set cid=:cid,filcount =:value,address=:to, status = :status,update_time=now(),errormsg=:error where id=:id `

	//sql1 := ` insert into message(cid,w_id, height, charge,status) values(:cid,:id,'0','0',3)  `

	sql = utils.SqlReplaceParames(sql, param)
	_, err := f.Engine.Exec(sql)
	if err != nil {
		log.Println("withdraw状态未改变：", param["id"])
		return err
	}
	//_,err = f.Engine.Exec(sql1)
	//if err!=nil{
	//	log.Println("withdraw消息未插入：",param["id"])
	//	return err
	//}
	return nil
}

//获取默认钱包地址
func (f *FinanceService) DefaultAddress(admin_id string) string {
	sql := ` select address,isdefault from wallet where admin_id='` + admin_id + `' and del_flag=0 order by isdefault desc,create_time asc limit 1 `

	res, err := f.Engine.QueryString(sql)
	if err != nil || len(res) == 0 {
		log.Println(err)
		return ""
	}

	return res[0]["address"]
}

//钱包余额
func (f *FinanceService) WalletBalance(admin_id string) (string, types.BigInt, error) {
	//获取默认钱包
	addrstr := f.DefaultAddress(admin_id)

	addrs, err := address.NewFromString(addrstr)
	if err != nil || addrstr == "" {
		return "", types.NewInt(0), err
	}
	balance, err := f.Rpc.WalletBalance(context.TODO(), addrs)
	return addrstr, balance, err
}

func (f *FinanceService) GetBatchList(param map[string]string) ([]map[string]string, error) {
	//sql := ` select * from batchnum where admin_id= `+param["admin_id"]
	sql := ` select operator,id,create_time,num,ifnull(surplus,0)surplus,count,b.admin_id,
       ifnull(c.amount,0)amount,ifnull(c.filcount,0)filcount,ifnull(c.charage,0)charage,ifnull(c.batch_count,0)batch_count,
       ifnull(d.charage1,0)charage1 from batchnum b
        left join (select count(1) surplus ,w.batch_count,w.admin_id from withdraw w where w.admin_id=:admin_id and status = 0 group by w.batch_count ) a
            on b.admin_id = a.admin_id and b.num = a.batch_count
        left join (select round(sum(w.amount), 4)   amount,
                        round(sum(w.filcount), 4) filcount,
                        round(sum(w.charge), 9)   charage,
                        w.batch_count
                 from withdraw w
                 where w.admin_id = :admin_id
                 group by w.batch_count)c on c.batch_count = b.num
        left join(select b.batch_count,sum(round(b.amount-b.filcount,4)) charage1
                    from withdraw b where b.admin_id=:admin_id and b.filcount<>0 group by b.batch_count )d
         on c.batch_count = d.batch_count
        where b.admin_id=:admin_id `
	if param["status"] == "0" {
		sql += " and ifnull(surplus,0) > 0"
	}
	if param["status"] != "0" && param["status"] != "" {
		sql += " and ifnull(surplus,0) = 0"
	}
	param["sort"] = "create_time"
	param["order"] = "desc"

	sql = utils.SqlReplaceParames(sql, param)
	//求总数
	var err error
	param["total"], err = utils.GetTotalCount(f.Engine, sql)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	sql += utils.LimitAndOrderBy(param)

	mp, err := f.Engine.QueryString(sql)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return mp, err
}

//var MSG = make(chan cid.Cid,2048)
//按批次转账
func (f *FinanceService) BatchSendsList(num string, user *entity.SysUser) (int, error) {
	adminid := strconv.FormatInt(user.AdminId, 10)
	//钱包地址from
	from, err := address.NewFromString(f.DefaultAddress(adminid))
	if err != nil {
		return 0, err
	}
	if num == "" {
		return 0, errors.New("参数不能为空！")
	}
	nums := strings.ReplaceAll(num, ",", "','")
	//查询数据库当前批次的用户
	sql := ` select id,address,amount from withdraw a
        where admin_id=` + adminid + `
		and batch_count in('` + nums + `') and status = 0  `

	users, err := f.Engine.QueryString(sql)
	if err != nil {
		return 0, err
	}
	config := utils.GetConfig(f.Engine, "charge", user)
	charge, err := strconv.ParseFloat(config, 64)
	if err != nil {
		charge = 0
	}
	n := 0

	param := make(map[string]string)
	for _, us := range users {
		param["status"] = "0"
		param["value"] = "0"
		//钱包地址to
		flage := true
		addrTo, err := address.NewFromString(us["address"])
		if err != nil {
			log.Println("地址格式异常：", err)
			//return 0,err
			param["cid"] = "地址格式异常"
			param["error"] = err.Error()
			flage = false
		}
		fil, _ := strconv.ParseFloat(us["amount"], 64)
		val, err := types.ParseFIL(strconv.FormatFloat(fil-fil*charge, 'f', -1, 64))
		if err != nil || fil < 0 {
			//return 0,fmt.Errorf("failed to parse amount: %w", err)
			log.Println("failed to parse amount：", err)
			param["cid"] = "金额格式不对"
			param["error"] = err.Error()
			flage = false
		}
		if flage {
			sm, err := f.Rpc.MpoolPushMessage(context.TODO(), &types.Message{
				To:    addrTo,
				From:  from,
				Value: types.BigInt(val),
				//GasLimit:   472835,
				//GasFeeCap:  abi.NewTokenAmount(101544),
				//GasPremium: abi.NewTokenAmount(100490),
			}, nil)
			if err != nil {
				param["cid"] = "钱包地址不正确或者余额不足："
				param["error"] = err.Error()
				n++
			} else {
				param["status"] = "3"
				param["cid"] = sm.Cid().String()
				param["error"] = ""
				param["value"] = strconv.FormatFloat(fil-fil*charge, 'f', -1, 64)
			}
		}
		param["id"] = us["id"]
		param["to"] = us["address"]
		err = f.UpdateTransferInfo(param)
		if err != nil {
			return 0, err
		}
		//MSG <- sm.Cid()
	}

	return n, nil
}

//按批次jujue
func (f *FinanceService) BatchRefuseList(ids, status string) error {
	if ids == "" || status == "" {
		return errors.New("参数不能为空！")
	}
	id := strings.ReplaceAll(ids, ",", "','")
	sql := ` update withdraw set status = ` + status + `,update_time=now() where batch_count in('` + id + `') and status = 0 `
	_, err := f.Engine.Exec(sql)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//获取配置
func (f *FinanceService) GetConfigList(admin_id string) map[string]string {

	sql := ` select c_name,c_value from sys_config where admin_id='` + admin_id + "'"
	res, err := f.Engine.QueryString(sql)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return utils.DataToKeyValue(res, "c_name", "c_value")
}

//设置配置
func (f *FinanceService) SetConfigList(admin_id string, param map[string]string) error {
	sql1 := `select count(1) count from sys_config where admin_id= '` + admin_id + "'"
	res, _ := f.Engine.QueryString(sql1)
	sees := f.Engine.NewSession()
	defer sees.Commit()
	if res[0]["count"] == "0" {
		sql2 := `insert into sys_config(c_name,c_value,admin_id)values`
		i := 0
		for k, v := range param {
			if i > 0 {
				sql2 += ","
			}
			sql2 += `('` + k + `','` + v + `','` + admin_id + `')`
			i++
		}
		_, err := f.Engine.Exec(sql2)
		if err != nil {
			sees.Rollback()
			return err
		}
	} else {
		for k, v := range param {
			sql := ` update sys_config set c_value='` + v + `' where c_name='` + k + `' and admin_id='` + admin_id + "'"
			_, err := f.Engine.Exec(sql)
			if err != nil {
				sees.Rollback()
				return err
			}
		}
	}

	return nil
}

//设置手续费
func (f *FinanceService) SetProcedureFee(cid, value string) error {
	sql := ` update withdraw set charge =:value,status = 1 where cid =:cid `
	param := make(map[string]string)
	param["cid"] = cid
	param["value"] = value
	sql = utils.SqlReplaceParames(sql, param)
	_, err := f.Engine.Exec(sql)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//获取正在转账批次，有当前批次就不进行转账
func (fi *FinanceService) GetSendsFlag(tokenStr, num string) (error, bool) {
	user, _ := utils.GetSubjectByTokenStr(tokenStr, fi.Client)
	if user.Sends != "" {
		sends := strings.Split(user.Sends, ",")
		nums := strings.Split(num, ",")
		mp := make(map[string]string)
		for _, v := range nums {
			mp[v] = ""
		}
		for i := 0; i < len(sends); i++ {
			if _, ok := mp[sends[i]]; ok {
				return errors.New("该批次正在转账，请不要重复提交！"), false
			}
		}
		user.Sends += ","
	}
	//锁定批次
	user.Sends += num

	userInfo, _ := json.Marshal(user)
	if err := fi.Client.Set(conf.TOKEN_PREFIX+tokenStr, string(userInfo), 2*time.Second*time.Duration(conf.TOKEN_EFFECT_TIME)).Err(); err != nil {
		return err, false
	}

	return nil, true
}

//获取删除转账完成的批次
func (fi *FinanceService) SetSendsByRedis(tokenStr, num string) error {
	user, _ := utils.GetSubjectByTokenStr(tokenStr, fi.Client)
	if user.Sends == "" {
		return nil
	}
	sends := strings.Split(user.Sends, ",")
	nums := strings.Split(num, ",")
	mp := make(map[string]string)
	for _, v := range nums {
		mp[v] = ""
	}
	str := ""
	for i := 0; i < len(sends); i++ {
		if _, ok := mp[sends[i]]; ok {
			continue
		}
		str += sends[i]
		if i < len(sends)-1 {
			str += ","
		}
	}
	//解锁
	user.Sends = str

	userInfo, _ := json.Marshal(user)
	if err := fi.Client.Set(conf.TOKEN_PREFIX+tokenStr, string(userInfo), 2*time.Second*time.Duration(conf.TOKEN_EFFECT_TIME)).Err(); err != nil {
		return err
	}

	return nil
}
