package handler

import (
	"net/http"

	"demacia/datacenter/parent/api/internal/logic"
	"demacia/datacenter/parent/api/internal/svc"
	"demacia/datacenter/parent/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func importParentsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UrlReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewImportParentsLogic(r.Context(), ctx)
		err := l.ImportParents(req)
		httpx.FormatResponse(nil, err, w)
	}
}
