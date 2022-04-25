package handler

import (
	"net/http"

	"demacia/datacenter/member/api/internal/logic"
	"demacia/datacenter/member/api/internal/svc"
	"demacia/datacenter/member/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func DeleteHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MemberIds
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewDeleteLogic(r.Context(), ctx)
		err := l.Delete(req)
		httpx.FormatResponse(nil, err, w)
	}
}
