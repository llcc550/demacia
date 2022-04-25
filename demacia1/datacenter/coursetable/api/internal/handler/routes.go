// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"demacia/datacenter/coursetable/api/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Log},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/api/coursetable/class_table",
					Handler: ClassCourseTableHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/coursetable/teacher_table",
					Handler: TeacherCourseTableHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/coursetable/position_table",
					Handler: PositionCourseTableHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/coursetable/deploy",
					Handler: CourseTableDeployHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/coursetable/deploy_save",
					Handler: CourseTableDeploySaveHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/coursetable/deploy_generate",
					Handler: GenerateCourseTableDeployHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/coursetable/info",
					Handler: CourseTableInfoHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/coursetable/add",
					Handler: CourseTableAddHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/coursetable/org_table",
					Handler: OrgCourseTableHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/coursetable/student_table",
					Handler: StudentCourseTableHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/coursetable/my_coursetable",
					Handler: MyCourseTableHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/coursetable/clear_deploy",
					Handler: ClearCourseDeployHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/class/list-teach",
					Handler: teachListHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/class/add-teach",
					Handler: addTeachHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/class/update-teach",
					Handler: updateTeachHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/class/delete-teach",
					Handler: deleteTeachHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}
