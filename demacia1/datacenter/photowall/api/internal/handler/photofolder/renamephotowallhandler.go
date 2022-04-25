package photofolder

import (
	"net/http"

	"demacia/datacenter/photowall/api/internal/logic/photofolder"
	"demacia/datacenter/photowall/api/internal/svc"
	"demacia/datacenter/photowall/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func RenamePhotowallHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdTitleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := photofolder.NewRenamePhotowallLogic(r.Context(), ctx)
		err := l.RenamePhotowall(req)
		httpx.FormatResponse(nil, err, w)
	}
}
