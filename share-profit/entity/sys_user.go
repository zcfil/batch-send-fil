package entity

import (
	"fmt"
	"time"
)

type JsonTime time.Time

func (jsonTime JsonTime)MarshalJSON()([]byte,error)  {

	jsonTimeStr := fmt.Sprintf(`"%s"`, time.Time(jsonTime).Format("2006-01-02 15:04:05"))

	fmt.Println(jsonTimeStr)

	return []byte(jsonTimeStr),nil

	//b := make([]byte, 0, len("2006-01-02 15:04:05")+2)
	//b = append(b, '"')
	//b = time.Time(jsonTime).AppendFormat(b, "2006-01-02 15:04:05")
	//b = append(b, '"')
	//return b, nil
}

func (jsonTime *JsonTime)UnmarshalJSON(data []byte)error  {

	t, err := time.Parse(`"2006-01-02 15:04:05"`, string(data))

	if err != nil {

		return err
	}

	*jsonTime = JsonTime(t)

	return nil
}

type Obj struct{

	From JsonTime `json:"from"`
}

type SysUser struct {
	Id         int64     `json:"userId,string" xorm:"pk comment('id') BIGINT(20)"`
	AdminId    int64     `json:"AdminId,string" xorm:"pk comment('id') BIGINT(20)"`
	Username   string    `json:"username" xorm:"comment('用户名') unique VARCHAR(50)"`
	Password   string    `json:"password" xorm:"comment('密码') VARCHAR(100)"`
	RealName   string    `json:"realName" xorm:"comment('姓名') VARCHAR(50)"`
	//Wallet    string    `json:"wallet" xorm:"comment('fil钱包地址') VARCHAR(200)"`
	HeadUrl    string    `json:"headUrl" xorm:"comment('头像') VARCHAR(200)"`
	Gender     int       `json:"gender" xorm:"comment('性别   0：男   1：女    2：保密') TINYINT(4)"`
	Email      string    `json:"email" xorm:"comment('邮箱') VARCHAR(100)"`
	Mobile     string    `json:"mobile" xorm:"comment('手机号') VARCHAR(20)"`
	DeptId     int64     `json:"deptId,string" xorm:"comment('部门ID') BIGINT(20)"`
	SuperAdmin int       `json:"superAdmin" xorm:"comment('超级管理员   0：否   1：是') TINYINT(3)"`
	Verified   int       `json:"verifield" xorm:"not null TINYINT(1)"`
	Status     int       `json:"status" xorm:"comment('状态  0：停用    1：正常') TINYINT(4)"`
	Remark     string    `json:"remark" xorm:"comment('备注') VARCHAR(200)"`
	DelFlag    int       `json:"delFlag" xorm:"comment('删除标识  0：未删除    1：删除') index TINYINT(4)"`
	Creator    int64     `json:"creator" xorm:"comment('创建者') BIGINT(20)"`
	CreateDate time.Time `json:"createDate" xorm:"comment('创建时间') index DATETIME"`
	Updater    int64     `json:"updater" xorm:"comment('更新者') BIGINT(20)"`
	UpdateDate time.Time `json:"updateDate" xorm:"comment('更新时间') DATETIME"`
	Verification string  `json:"verification" xorm:"comment('谷歌验证') VARCHAR(255)""` //google验证
	Sends 		string  `json:"sends"` //正在进行转账交易的批次
}

type UserDept struct {
	SysUser `xorm:"extends"json:"sys_user"`
	SysDept `xorm:"extends" json:"sys_Dept"`
}


/**
func (userDept *UserDept)MarshalJSON() ([]byte, error) {

	//buf := new(bytes.Buffer)
	//var id int64 = 12457895478541
	//err := binary.Write(buf, binary.LittleEndian, &res)
	//if err != nil {
	//	fmt.Println("binary.Write failed:", err)
	//}
	//return []byte(strconv.FormatInt(123456789,10)),nil

	//buf := bytes.NewBuffer(nil)
	//enc := gob.NewEncoder(buf)
	//err := enc.Encode(res)
	//if err != nil {
	//	return nil, err
	//}
	//log.Println(buf.Bytes())

	//return buf.Bytes(), nil

	formatString := userDept.SysUser.CreateDate.Format("2006-01-02 15:04:05")
	userDept.SysUser.CreateDate

	return nil, nil

}
*/

