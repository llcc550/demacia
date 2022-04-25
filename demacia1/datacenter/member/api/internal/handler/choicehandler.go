package handler

import (
	"net/http"

	"demacia/datacenter/member/api/internal/logic"
	"demacia/datacenter/member/api/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func ChoiceHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewChoiceLogic(r.Context(), ctx)
		resp, err := l.Choice()
		httpx.FormatResponse(resp, err, w)
	}
}
