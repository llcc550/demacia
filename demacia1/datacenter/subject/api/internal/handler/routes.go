// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"demacia/datacenter/subject/api/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Log},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/api/subject/list",
					Handler: SubjectListHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/subject/insert",
					Handler: AddSubjectHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/subject/Rename",
					Handler: RenameHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/subject/GradeManage",
					Handler: GradeManageHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/subject/delete",
					Handler: DeletedSubjectHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/subject/TeacherManage",
					Handler: TeacherManageHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/subject/getSubjectTeacher",
					Handler: GetSubjectTeacherHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}