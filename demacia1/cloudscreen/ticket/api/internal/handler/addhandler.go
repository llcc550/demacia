package handler

import (
	"net/http"

	"demacia/cloudscreen/ticket/api/internal/logic"
	"demacia/cloudscreen/ticket/api/internal/svc"
	"demacia/cloudscreen/ticket/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func addHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewAddLogic(r.Context(), ctx)
		err := l.Add(req)
		httpx.FormatResponse(nil, err, w)
	}
}
