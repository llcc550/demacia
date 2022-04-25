package handler

import (
	"net/http"

	"demacia/datacenter/organization/api/internal/logic"
	"demacia/datacenter/organization/api/internal/svc"
	"demacia/datacenter/organization/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func listHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrgListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewListLogic(r.Context(), ctx)
		resp, err := l.List(req)
		httpx.FormatResponse(resp, err, w)
	}
}
