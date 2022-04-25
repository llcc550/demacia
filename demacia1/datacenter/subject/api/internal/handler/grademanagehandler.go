package handler

import (
	"demacia/datacenter/subject/api/internal/logic"
	"demacia/datacenter/subject/api/internal/svc"
	"demacia/datacenter/subject/api/internal/types"
	"net/http"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func GradeManageHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GradeManageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewGradeManageLogic(r.Context(), ctx)
		err := l.GradeManage(req)
		httpx.FormatResponse(nil, err, w)
	}
}
