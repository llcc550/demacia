package handler

import (
	"net/http"

	"demacia/datacenter/class/api/internal/logic"
	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func addTeachHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddTeachReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewAddTeachLogic(r.Context(), ctx)
		resp, err := l.AddTeach(req)
		httpx.FormatResponse(resp, err, w)
	}
}
