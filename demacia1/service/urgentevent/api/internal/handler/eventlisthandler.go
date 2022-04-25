package handler

import (
	"net/http"

	"demacia/service/urgentevent/api/internal/logic"
	"demacia/service/urgentevent/api/internal/svc"
	"demacia/service/urgentevent/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func EventListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewEventListLogic(r.Context(), ctx)
		resp, err := l.EventList(req)
		httpx.FormatResponse(resp, err, w)
	}
}
