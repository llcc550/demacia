package handler

import (
	"net/http"

	"demacia/datacenter/organization/api/internal/logic"
	"demacia/datacenter/organization/api/internal/svc"
	"demacia/datacenter/organization/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func addHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrgAddReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewAddLogic(r.Context(), ctx)
		resp, err := l.Add(req)
		httpx.FormatResponse(resp, err, w)
	}
}
