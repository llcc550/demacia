package handler

import (
	"net/http"

	"demacia/service/auth/api/internal/logic"
	"demacia/service/auth/api/internal/svc"
	"demacia/service/auth/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func userLoginHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewUserLoginLogic(r.Context(), ctx)
		resp, err := l.UserLogin(req)
		httpx.FormatResponse(resp, err, w)
	}
}
