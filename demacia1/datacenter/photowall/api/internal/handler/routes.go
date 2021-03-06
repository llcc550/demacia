// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	photofolder "demacia/datacenter/photowall/api/internal/handler/photofolder"
	photos "demacia/datacenter/photowall/api/internal/handler/photos"
	"demacia/datacenter/photowall/api/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/photofolder/list",
				Handler: photofolder.ListPhotowallHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/photofolder/Insert",
				Handler: photofolder.InsertPhotowallHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/photofolder/rename",
				Handler: photofolder.RenamePhotowallHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/photofolder/delete",
				Handler: photofolder.DeletedPhotowallHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/photos/photolist",
				Handler: photos.ListPhotoHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/photos/Rename",
				Handler: photos.RenameHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/photos/editscreensaver",
				Handler: photos.EditscreensaverHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/photos/editlockscreen",
				Handler: photos.EditlockscreenHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/photos/edittopping",
				Handler: photos.EdittoppingHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/photos/editpublish",
				Handler: photos.EditpublishHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/photos/delete",
				Handler: photos.DeleteHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}
