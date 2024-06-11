package service

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/go-redis/redis/v7"
	"github.com/go-xorm/xorm"
	"go.uber.org/dig"
	"log"
	"share-profit/pkg"
	"share-profit/rpc"
	"share-profit/utils"
	"strconv"
	"strings"
)

type YungoService struct{
	dig.In
	*redis.Client
	*rpc.YungoRpc
	*xorm.Engine
}
//判断助记词是否存在
func (s *YungoService)IsMnemonic(mnemonic,admin_id string)(bool) {
	//判断该用户是否拥有该钱包
	sql :=  `select 1 from wallet where mnemonic =? and admin_id=? and del_flag=0 `
	res,err  :=  s.Engine.QueryString(sql,mnemonic,admin_id)
	if err!=nil||len(res) == 0{
		return false
	}

	return true
}
//获取钱包列表
func (s *YungoService)GetWalletList(param map[string]string) ([]map[string]string,error) {

	sql := ` select address,create_time,
			if(address=(select address from wallet where admin_id=:admin_id and del_flag=0 order by isdefault desc,create_time asc limit 1),1,0) isdefault
    			from wallet where admin_id=:admin_id and del_flag=0 `
	param["sort"] = "create_time"
	param["order"] = "desc"

	sql = utils.SqlReplaceParames(sql,param)

	//求总数
	var err error
	param["total"],err = utils.GetTotalCount(s.Engine,sql)
	if err!=nil{
		return nil,err
	}

	sql += utils.LimitAndOrderBy(param)

	mp,err := s.Engine.QueryString(sql)
	if err!=nil{
		return nil,err
	}
	for i,v := range mp{
		addr ,_ := address.NewFromString(v["address"])
		balance,err := s.Rpc.WalletBalance(context.TODO(),addr)
		if err!=nil{
			mp[i]["balance"] = "0"
			log.Println("获取余额错误：",addr,err)
			continue
		}
		mp[i]["balance"] = strconv.FormatFloat(utils.NanoOrAttoToFIL(balance.String(),utils.AttoFIL),'f',-1,64)
	}

	return mp, err
}

//创建 钱包
func (s *YungoService)WalletNew(param map[string]string)(string,error) {

	//创建钱包地址
	key,err := pkg.CreateKey(param["mnemonic"] )
	if err != nil {
		return "",err
	}
	param["address"] = key.Address.String()
	//保存钱包地址
	sess := s.Engine.NewSession()
	sess.Begin()
	defer func() {
		if err!=nil{
			sess.Rollback()
		}else{
			sess.Commit()
		}
	}()
	//因为本地服务器lotus出问题了，所以得远程调用，组装参数
	var kk types.KeyInfo
	bb,_ := json.Marshal(key.KeyInfo)
	err = json.Unmarshal(bb,kk)
	sql  := `insert into wallet (admin_id,address,mnemonic)value(:admin_id,:address,:mnemonic_ase) `
	sql = utils.SqlReplaceParames(sql,param)
	_,err = s.Engine.Exec(sql)
	if err != nil {
		return "",err
	}
	if addr,err := s.Rpc.WalletImport(context.Background(),&kk);err!=nil{
		log.Println(addr,err)
	}
	//err = key.KeyInfo.Put(pkg.KNamePrefix +key.Address.String())
	//if err != nil {
	//	return "",err
	//}

	return key.Address.String(),nil
}
//删除钱包
func (s *YungoService)DeleteWallet(param map[string]string)(error) {
	sql := `update wallet set del_flag = 1 where admin_id = :admin_id and address = :wallet `
	sql = utils.SqlReplaceParames(sql,param)
	_,err := s.Engine.Exec(sql)
	return  err
}
//设置默认钱包
func (s *YungoService)SetDefault(param map[string]string)(error) {
	sql := `update wallet a, (select max(b.isdefault)+1 isdefault
                from wallet b where b.admin_id=:admin_id ) b
            set a.isdefault=b.isdefault where address=:wallet and admin_id=:admin_id`
	sql = utils.SqlReplaceParames(sql,param)
	_,err := s.Engine.Exec(sql)
	return  err
}
//导入钱包
func (s *YungoService)ImprotWalletAddress(param map[string]string)(string,error) {

	b ,err := hex.DecodeString(param["private_key"])
	var keyInfo pkg.KeyInfo
	sss := strings.TrimSpace(string(utils.AesDecrypt(b)))
	str,_ := hex.DecodeString(sss)
	err = json.Unmarshal(str,&keyInfo)


	//因为本地服务器lotus出问题了，所以得远程调用，组装参数
	var kk types.KeyInfo
	err = json.Unmarshal(str,&kk)
	//fmt.Println("云构解密:",str,keyInfo)
	if err != nil {
		return "",errors.New("导入私钥有误！"+err.Error())
	}
	//key,err := pkg.NewKey(keyInfo)
	//if err != nil {
	//	return "",err
	//}
	//保存钱包地址
	sess := s.Engine.NewSession()
	sess.Begin()
	defer func() {
		if err!=nil{
			sess.Rollback()
		}else{
			sess.Commit()
		}
	}()
	addr,err := s.Rpc.WalletImport(context.TODO(),&kk)
	if err!=nil{
		log.Println(addr,err)
		if !strings.Contains(err.Error(),"key already exists"){
			return "",err
		}
		ads := strings.Split(err.Error(),"'")
		ad := ""
		for _,v := range ads{
			if strings.Contains(v,"wallet-"){
				ad = strings.Replace(v,"wallet-","",1)
			}
		}
		param["address"] = ad
		param["wallet"] = ad
	}else{
		param["address"] = addr.String()
		param["wallet"] = addr.String()
	}

	if s.IsPrivateKey(param){
		return "",errors.New("该地址已导入！")
	}

	sql  := `insert into wallet (admin_id,address,mnemonic)value(:admin_id,:address,'') `
	sql = utils.SqlReplaceParames(sql,param)
	_,err = s.Engine.Exec(sql)
	if err != nil {
		return "",err
	}

	//err = key.KeyInfo.Put(pkg.KNamePrefix +addr.String())
	//if err != nil {
	//	log.Println(err)
	//	//return "",nil
	//}

	return addr.String(),nil
}
//判断钱包是否存在私钥
func (s *YungoService)IsPrivateKey(param map[string]string)(bool) {
	//判断该用户是否拥有该钱包
	sql :=  `select 1 from wallet where admin_id=:admin_id and address =:wallet and del_flag=0 `
	sql = utils.SqlReplaceParames(sql,param)
	res,err  :=  s.Engine.QueryString(sql)
	if err!=nil||len(res) == 0{
		fmt.Println(err,res)
		return false
	}

	return true
}
//判断钱包是否还有人在用该钱包
func (s *YungoService)IsUseWallet(param map[string]string)(bool) {
	//判断该用户是否拥有该钱包
	sql :=  `select 1 from wallet where address =:wallet and del_flag=0 `
	sql = utils.SqlReplaceParames(sql,param)
	res,err  :=  s.Engine.QueryString(sql)
	if err!=nil||len(res) <= 1{
		return false
	}
	return true
}
//获取默认钱包地址
func (s *YungoService)DefaultAddress(param map[string]string)(string) {
	sql :=  ` select address,isdefault from wallet where admin_id=:admin_id order by isdefault desc,create_time asc limit 1 `
	sql = utils.SqlReplaceParames(sql,param)
	res,err  :=  s.Engine.QueryString(sql)
	if err!=nil||len(res) == 0{
		log.Println(err)
		return ""
	}

	return res[0]["address"]
}

//导出助记词
func (s *YungoService)ExportMnemonicByAddres(param map[string]string)(string,error) {
	sql :=  `select mnemonic from wallet where admin_id=:admin_id and address =:wallet `
	sql = utils.SqlReplaceParames(sql,param)
	res,err  :=  s.Engine.QueryString(sql)
	if err!=nil||len(res) == 0{
		return "",errors.New("钱包地址不存在! ")
	}

	return res[0]["mnemonic"],nil
}

//导入钱包 助记词
func (s *YungoService)ImprotMnemonicAddress(param map[string]string)(string,error) {
	//创建钱包地址
	key,err := pkg.CreateKey(param["mnemonic"] )
	if err != nil {
		return "",err
	}
	param["address"] = key.Address.String()
	param["wallet"] = key.Address.String()
	if s.IsPrivateKey(param){
		return "",errors.New("该地址已存在！")
	}
	//保存钱包地址
	sess := s.Engine.NewSession()
	sess.Begin()
	defer func() {
		if err!=nil{
			sess.Rollback()
		}else{
			sess.Commit()
		}
	}()
	sql  := `insert into wallet (admin_id,address,mnemonic)value(:admin_id,:address,:mnemonic_ase) `
	sql = utils.SqlReplaceParames(sql,param)
	_,err = s.Engine.Exec(sql)
	if err != nil {
		return "",err
	}
	//因为本地服务器lotus出问题了，所以得远程调用，组装参数
	var kk types.KeyInfo
	bb,_ := json.Marshal(key.KeyInfo)
	err = json.Unmarshal(bb,kk)
	if addr,err := s.Rpc.WalletImport(context.Background(),&kk);err!=nil{
		log.Println(addr,err)
	}
	//err = key.KeyInfo.Put(pkg.KNamePrefix +key.Address.String())
	//if err != nil {
	//	log.Println(err)
	//}

	return key.Address.String(),nil
}
func (s *YungoService)CheckMnemonic(mnemonic string)bool{
	ms := strings.Split(mnemonic," ")
	if len(ms)!= pkg.MnemonicCount{
		return false
	}
	for _,v := range ms{
		flage := false
		for _,val := range pkg.English{
			if v == val{
				flage = true
				break
			}
		}
		if !flage{
			return false
		}
	}
	return true
}
