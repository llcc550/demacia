package handler

import (
	"net/http"

	"demacia/datacenter/class/api/internal/logic"
	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func deleteGradeByIdHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Id
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewDeleteGradeByIdLogic(r.Context(), ctx)
		err := l.DeleteGradeById(req)
		httpx.FormatResponse(nil, err, w)
	}
}
