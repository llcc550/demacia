package photos

import (
	"net/http"

	"demacia/datacenter/photowall/api/internal/logic/photos"
	"demacia/datacenter/photowall/api/internal/svc"
	"demacia/datacenter/photowall/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func EditlockscreenHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EditReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := photos.NewEditlockscreenLogic(r.Context(), ctx)
		err := l.Editlockscreen(req)
		httpx.FormatResponse(nil, err, w)
	}
}
