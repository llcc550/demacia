// Code generated by goctl. DO NOT EDIT.
package types

type UrlReq struct {
	WebsocketUuid string `json:"websocket_uuid"`
	Url           string `json:"url"`
}

type InsertRequest struct {
	ParentId   int64     `json:"parent_id,optional"`
	Student    []Student `json:"student"`
	ParentName string    `json:"parent_name"`
	Moblie     string    `json:"mobile"`
	IdNumber   string    `json:"id_number,optional"`
	Address    string    `json:"address,optional"`
	Face       string    `json:"face,optional"`
}

type Student struct {
	StudentParentId int64 `json:"student_parent_id,optional"`
	ParentId        int64 `json:"parent_id,optional"`
	StudentId       int64 `json:"student_id"`
	ClassId         int64 `json:"class_id"`
	Relation        int8  `json:"relation"`
}

type IdRequest struct {
	ParentId int64 `json:"parent_id"`
}

type IdsRequest struct {
	ParentIds []int64 `json:"parent_ids"`
}

type List struct {
	ParentId    int64         `json:"parent_id"`
	ParentName  string        `json:"parent_name"`
	Mobile      string        `json:"mobile"`
	StudentInfo []StudentInfo `json:"student_info"`
	FaceStatus  int8          `json:"face_status"`
}

type StudentInfo struct {
	StudentId       int64  `json:"student_id"`
	StudentName     string `json:"student_name"`
	StudentUserName string `json:"student_user_name"`
	ClassName       string `json:"class_name"`
	Relation        int8   `json:"relation"`
}

type ParentList struct {
	List  []List `json:"list"`
	Count int    `json:"count"`
}

type ListConditionRequest struct {
	ClassId     int64  `json:"class_id,optional"`
	StudentName string `json:"student_name,optional"`
	ParentName  string `json:"parent_name,optional"`
	FaceStatus  int8   `json:"face_status,optional"`
	Page        int    `json:"page,optional"`
	Limit       int    `json:"limit,optional"`
}

type DetailResponse struct {
	ParentId    int64         `json:"parent_id"`
	StudentInfo []StudentInfo `json:"student_info"`
	ParentName  string        `json:"parent_name"`
	Moblie      string        `json:"mobile"`
	IdNumber    string        `json:"id_number"`
	Address     string        `json:"address"`
	Face        string        `json:"face"`
}
