// Code generated by goctl. DO NOT EDIT.
package types

type Area struct {
	Id   int64  `json:"id"`
	Pid  int64  `json:"pid"`
	Name string `json:"name"`
}

type AreaReq struct {
	Pid int64 `json:"pid,optional"`
}

type AreaResp struct {
	List []*Area `json:"list"`
}

type Ethnic struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type EthnicResp struct {
	List []*Ethnic `json:"list"`
}