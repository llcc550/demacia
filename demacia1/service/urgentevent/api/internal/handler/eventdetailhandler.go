package handler

import (
	"net/http"

	"demacia/service/urgentevent/api/internal/logic"
	"demacia/service/urgentevent/api/internal/svc"
	"demacia/service/urgentevent/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func EventDetailHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EventIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewEventDetailLogic(r.Context(), ctx)
		resp, err := l.EventDetail(req)
		httpx.FormatResponse(resp, err, w)
	}
}
