package handler

import (
	"net/http"

	"demacia/datacenter/device/api/internal/logic"
	"demacia/datacenter/device/api/internal/svc"
	"demacia/datacenter/device/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func groupDetailHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewGroupDetailLogic(r.Context(), ctx)
		resp, err := l.GroupDetail(req)
		httpx.FormatResponse(resp, err, w)
	}
}
