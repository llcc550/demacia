package handler

import (
	"net/http"

	"demacia/datacenter/organization/api/internal/logic"
	"demacia/datacenter/organization/api/internal/svc"
	"demacia/datacenter/organization/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func delHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewDelLogic(r.Context(), ctx)
		err := l.Del(req)
		httpx.FormatResponse(nil, err, w)
	}
}
