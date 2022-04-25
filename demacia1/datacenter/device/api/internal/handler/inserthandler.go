package handler

import (
	"net/http"

	"demacia/datacenter/device/api/internal/logic"
	"demacia/datacenter/device/api/internal/svc"
	"demacia/datacenter/device/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func insertHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.InsertReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewInsertLogic(r.Context(), ctx)
		err := l.Insert(req)
		httpx.FormatResponse(nil, err, w)
	}
}
