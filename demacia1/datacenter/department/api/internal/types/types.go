// Code generated by goctl. DO NOT EDIT.
package types

type ListReq struct {
	Title string `json:"title,optional"`
}

type ListInfo struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Sort        int64  `json:"sort"`
	MemberCount int64  `json:"member_count"`
}

type ListResp struct {
	List []*ListInfo `json:"list"`
}

type InsertReq struct {
	Title string `json:"title"`
}

type UpdateReq struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
}

type DelReq struct {
	Ids []int64 `json:"ids"`
}

type SortInfo struct {
	Id   int64 `json:"id"`
	Sort int64 `json:"sort"`
}

type SortReq struct {
	List []*SortInfo `json:"list"`
}

type MemberListReq struct {
	DepartmentId int64 `json:"department_id"`
	Page         int   `json:"page"`
	Limit        int   `json:"limit"`
}

type MemberInfo struct {
	MemberId int64  `json:"member_id"`
	TrueName string `json:"true_name"`
	Mobile   string `json:"mobile"`
}

type MemberListResp struct {
	Count int64         `json:"count"`
	List  []*MemberInfo `json:"list"`
}

type MemberIdsReq struct {
	DepartmentId int64   `json:"department_id"`
	MemberIds    []int64 `json:"member_ids"`
}
