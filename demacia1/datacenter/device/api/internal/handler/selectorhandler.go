package handler

import (
	"net/http"

	"demacia/datacenter/device/api/internal/logic"
	"demacia/datacenter/device/api/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func selectorHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewSelectorLogic(r.Context(), ctx)
		resp, err := l.Selector()
		httpx.FormatResponse(resp, err, w)
	}
}
