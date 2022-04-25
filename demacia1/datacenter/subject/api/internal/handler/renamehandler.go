package handler

import (
	"demacia/datacenter/subject/api/internal/logic"
	"demacia/datacenter/subject/api/internal/svc"
	"demacia/datacenter/subject/api/internal/types"
	"net/http"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func RenameHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Info
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewRenameLogic(r.Context(), ctx)
		err := l.Rename(req)
		httpx.FormatResponse(nil, err, w)
	}
}
