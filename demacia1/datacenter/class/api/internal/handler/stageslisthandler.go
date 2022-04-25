package handler

import (
	"net/http"

	"demacia/datacenter/class/api/internal/logic"
	"demacia/datacenter/class/api/internal/svc"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func stagesListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewStagesListLogic(r.Context(), ctx)
		resp, err := l.StagesList()
		httpx.FormatResponse(resp, err, w)
	}
}
