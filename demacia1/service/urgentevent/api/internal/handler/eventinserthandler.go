package handler

import (
	"net/http"

	"demacia/service/urgentevent/api/internal/logic"
	"demacia/service/urgentevent/api/internal/svc"
	"demacia/service/urgentevent/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func EventInsertHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EventInsertReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewEventInsertLogic(r.Context(), ctx)
		err := l.EventInsert(req)
		httpx.FormatResponse(nil, err, w)
	}
}
