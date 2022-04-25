package handler

import (
	"net/http"

	"demacia/service/event/api/internal/logic"
	"demacia/service/event/api/internal/svc"
	"demacia/service/event/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func CategoryUpdateHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CategoryUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewCategoryUpdateLogic(r.Context(), ctx)
		err := l.CategoryUpdate(req)
		httpx.FormatResponse(nil, err, w)
	}
}
