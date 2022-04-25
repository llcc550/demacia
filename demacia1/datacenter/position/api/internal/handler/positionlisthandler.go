package handler

import (
	"demacia/common/errlist"
	"net/http"

	"demacia/datacenter/position/api/internal/logic"
	"demacia/datacenter/position/api/internal/svc"
	"demacia/datacenter/position/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func PositionListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, errlist.InvalidParam, w)
			return
		}

		l := logic.NewPositionListLogic(r.Context(), ctx)
		resp, err := l.PositionList(req)
		httpx.FormatResponse(resp, err, w)
	}
}
