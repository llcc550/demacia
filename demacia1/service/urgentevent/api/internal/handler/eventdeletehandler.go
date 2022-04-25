package handler

import (
	"net/http"

	"demacia/service/urgentevent/api/internal/logic"
	"demacia/service/urgentevent/api/internal/svc"
	"demacia/service/urgentevent/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func EventDeleteHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EventIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewEventDeleteLogic(r.Context(), ctx)
		err := l.EventDelete(req)
		httpx.FormatResponse(nil, err, w)
	}
}
