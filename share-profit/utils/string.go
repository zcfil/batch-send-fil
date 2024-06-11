package utils

import (
	"strconv"
	"strings"
)

func LimitAndOrderBy(param map[string]string)string{
	str := ""
	//排序
	if param["sort"]!=""{
		str += ` order by `+param["sort"]
		if param["order"]!=""{
			str +=" "+ param["order"]
		}
	}
	if param["isexp"]==""{
		param["isexp"]="0"
	}
	if param["isexp"]!="1"{
		//分页
		if param["pageNum"]!=""&&param["pageSize"]!=""{
			pageNum,_ := strconv.Atoi(param["pageNum"])
			pageSize,_:= strconv.Atoi(param["pageSize"])
			if pageNum!=0 && pageSize!=0{
				str += ` limit `+strconv.Itoa((pageNum-1)*pageSize)+`,`+param["pageSize"]
			}
		}

	}

	return str
}

//将sql中的占位符':'替换成map中的参数
func SqlReplaceParames(sql string,param map[string]string)(string){
	fa := false
	start := 0
	sqlstr := sql
	fl := true
	for i,v := range sql{
		if v==':'{
			start = i+1
			fa = true
		}
		if (v == '\n'|| v=='\t'||v==' '||v==','||v==')'||v=='%'||v=='"'||v=='='||len(sql)-1==i)&&fa&&fl{
			field := sql[start:i]
			//最后一个
			if len(sql)-1==i&&v!=' '&&v != '\n'&& v!='\t'&&v!=')'{
				field = sql[start:i+1]
			}
			if param[field]!=""{
				if sql[start-3]=='%'{
					sqlstr = strings.Replace(sqlstr,"%%:"+field+"%%",`'%%`+param[field]+`%%'`,-1)
				}else if sql[start-2]=='%' {
					sqlstr = strings.Replace(sqlstr,"%:"+field+"%",`'%`+param[field]+`%'`,-1)
				}else{
					sqlstr = strings.Replace(sqlstr,":"+field,`'`+param[field]+`'`,-1)
				}
				fa = false
			}else{
				if _,ok := param[field];ok{
					sqlstr = strings.Replace(sqlstr,":"+field,`'`+param[field]+`'`,-1)
					fa = false
					continue
				}
				if sql[i-1]=='\''||sql[i-1]=='"'{
					fa = false
					continue
				}

				sqlstr = sql[start-5:i+10] + " 参数不存在!"+strconv.Itoa(i)
				return sqlstr
			}
		}

	}
	return sqlstr
}
func CheckPrefix(content string,prefix string)bool{
	if len(content) < len(prefix){
		return false
	}
	if content[0:len(prefix)]!=prefix{
		return false
	}
	return true
}
