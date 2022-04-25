// Code generated by goctl. DO NOT EDIT.
package types

import "github.com/tealeg/xlsx"

type UrlReq struct {
	WebsocketUuid string `json:"websocket_uuid"`
	Url           string `json:"url"`
}

type InsertRequest struct {
	StudentId   int64    `json:"student_id,optional"`
	ClassId     int64    `json:"class_id"`
	UserName    string   `json:"user_name"`
	StudentName string   `json:"student_name"`
	Sex         int8     `json:"sex"`
	CardNumber  []string `json:"card_number,optional"`
	IdNumber    string   `json:"id_number,optional"`
	Avatar      string   `json:"avatar,optional"`
	Face        string   `json:"face,optional"`
}

type IdsRequest struct {
	StudentIds []int64 `json:"student_ids"`
}

type ListConditionRequest struct {
	OrgId       int64  `json:"org_id,optional"`
	StageId     int64  `json:"stage_id,optional"`
	ClassId     int64  `json:"class_id,optional"`
	GradeId     int64  `json:"grade_id,optional"`
	StudentName string `json:"student_name,optional"`
	FaceStatus  int8   `json:"face_status,optional"`
	Page        int    `json:"page,optional"`
	Limit       int    `json:"limit,optional"`
}

type List struct {
	StudentId   int64  `json:"student_id"`
	StudentName string `json:"student_name"`
	ClassName   string `json:"class_name"`
	UserName    string `json:"user_name"`
	Sex         int8   `json:"sex"`
	Face        string `json:"face"`
}

type ListResponse struct {
	List  []*List `json:"list"`
	Count int     `json:"count"`
}

type IdRequest struct {
	StudentId int64 `json:"student_id"`
}

type StudentDetail struct {
	StudentId   int64    `json:"student_id"`
	ClassId     int64    `json:"class_id"`
	ClassName   string   `json:"class_name"`
	UserName    string   `json:"user_name"`
	StudentName string   `json:"student_name"`
	Sex         int8     `json:"sex"`
	CardNumber  []string `json:"card_number"`
	IdNumber    string   `json:"id_number"`
	Avatar      string   `json:"avatar"`
	Face        string   `json:"face"`
}

type ClassIdRequest struct {
	ClassId int64 `json:"class_id"`
}

type TemplateFile struct {
	File *xlsx.File
	Name string `json:"name"`
}