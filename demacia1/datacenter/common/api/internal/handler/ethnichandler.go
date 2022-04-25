package handler

import (
	"net/http"

	"demacia/datacenter/common/api/internal/logic"
	"demacia/datacenter/common/api/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func ethnicHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewEthnicLogic(r.Context(), ctx)
		resp, err := l.Ethnic()
		httpx.FormatResponse(resp, err, w)
	}
}
