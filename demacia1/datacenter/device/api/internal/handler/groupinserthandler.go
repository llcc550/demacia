package handler

import (
	"net/http"

	"demacia/datacenter/device/api/internal/logic"
	"demacia/datacenter/device/api/internal/svc"
	"demacia/datacenter/device/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func groupInsertHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupInsertReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewGroupInsertLogic(r.Context(), ctx)
		err := l.GroupInsert(req)
		httpx.FormatResponse(nil, err, w)
	}
}
