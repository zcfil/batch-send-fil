package utils

import (
	"encoding/json"
	"strconv"
)

//Pagination 分页数据
type Pagination struct {
	PageNo   int         `json:"pageNo"`
	PageSize int         `json:"pageSize"`
	List     interface{} `json:"list"`
	Total    int         `json:"total"`
}
type PageData struct {
	PageNo   int         `json:"pageNo"`
	PageSize int         `json:"pageSize"`
	List     interface{} `json:"list"`
	Total    int         `json:"total"`
}
func NewPagination(pageNo, pageSize, total int, list interface{}) *Pagination {

	return &Pagination{
		PageNo:   pageNo,
		PageSize: pageSize,
		List:     list,
		Total:    total,
	}
}

//Pagination
func NewPageData(param map[string]string, list interface{}) *PageData {
	no,_ := strconv.Atoi(param["pageNum"])
	size,_ := strconv.Atoi(param["pageSize"])
	to,_ := strconv.Atoi(param["total"])
	return &PageData{
		PageNo:   no,
		PageSize: size,
		List:     list,
		Total:    to,
	}
}

func (p *Pagination)String() string {
	bytes, _ := json.Marshal(p)
	return string(bytes)
}

