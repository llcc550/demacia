package photofolder

import (
	"net/http"

	"demacia/datacenter/photowall/api/internal/logic/photofolder"
	"demacia/datacenter/photowall/api/internal/svc"
	"demacia/datacenter/photowall/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func ListPhotowallHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := photofolder.NewListPhotowallLogic(r.Context(), ctx)
		resp, err := l.ListPhotowall(req)
		httpx.FormatResponse(resp, err, w)
	}
}
