package handler

import (
	"net/http"

	"demacia/datacenter/coursetable/api/internal/logic"
	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func teachListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Id
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewTeachListLogic(r.Context(), ctx)
		resp, err := l.TeachList(req)
		httpx.FormatResponse(resp, err, w)
	}
}
