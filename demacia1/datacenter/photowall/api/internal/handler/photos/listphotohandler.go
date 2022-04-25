package photos

import (
	"net/http"

	"demacia/datacenter/photowall/api/internal/logic/photos"
	"demacia/datacenter/photowall/api/internal/svc"
	"demacia/datacenter/photowall/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func ListPhotoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListPhotoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := photos.NewListPhotoLogic(r.Context(), ctx)
		resp, err := l.ListPhoto(req)
		httpx.FormatResponse(resp, err, w)
	}
}
