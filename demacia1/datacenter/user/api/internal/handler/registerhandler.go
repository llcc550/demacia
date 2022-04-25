package handler

import (
	"net/http"

	"demacia/datacenter/user/api/internal/logic"
	"demacia/datacenter/user/api/internal/svc"
	"demacia/datacenter/user/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func registerHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewRegisterLogic(r.Context(), ctx)
		err := l.Register(req)
		httpx.FormatResponse(nil, err, w)
	}
}
