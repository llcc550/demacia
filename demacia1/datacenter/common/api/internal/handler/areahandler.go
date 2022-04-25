package handler

import (
	"net/http"

	"demacia/datacenter/common/api/internal/logic"
	"demacia/datacenter/common/api/internal/svc"
	"demacia/datacenter/common/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func areaHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AreaReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewAreaLogic(r.Context(), ctx)
		resp, err := l.Area(req)
		httpx.FormatResponse(resp, err, w)
	}
}
