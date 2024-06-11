package result


type Result struct{

	Code int `json:"code"`
	Data interface{} `json:"data"`
	Msg string	`json:"msg"`
}

//Ok 请求成功的result
func Ok(data interface{}) Result {

	return Result{
		Code: 0,
		Data: data,
		Msg:  "success",
	}
}

//Fail 请求失败的result
func Fail(err error) Result {

	return Result{
		Code: -1,
		Data: nil,
		Msg:  err.Error(),
	}
}
