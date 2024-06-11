package models

//谷歌验证器
type GoogleAuth struct {
	GoogleSecret  string `json:"googleSecret,omitempty"`
	GoogleAccount string `json:"googleAccount,omitempty"`
	URL           string `json:"url,omitempty"`
	Code          string `json:"code,omitempty"`
	IsBind        bool   `json:"isBind,omitempty"`
}