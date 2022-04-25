package handler

import (
	"net/http"

	"demacia/datacenter/student/api/internal/logic"
	"demacia/datacenter/student/api/internal/svc"
	"demacia/datacenter/student/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func studentsDetailHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewStudentsDetailLogic(r.Context(), ctx)
		resp, err := l.StudentsDetail(req)
		httpx.FormatResponse(resp, err, w)
	}
}
