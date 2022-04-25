package handler

import (
	"net/http"

	"demacia/datacenter/member/api/internal/logic"
	"demacia/datacenter/member/api/internal/svc"
	"demacia/datacenter/member/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func DetailHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MemberId
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewDetailLogic(r.Context(), ctx)
		resp, err := l.Detail(req)
		httpx.FormatResponse(resp, err, w)
	}
}
