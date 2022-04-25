package handler

import (
	"net/http"

	"demacia/datacenter/student/api/internal/logic"
	"demacia/datacenter/student/api/internal/svc"
	"demacia/datacenter/student/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func classStudentsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ClassIdRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewClassStudentsLogic(r.Context(), ctx)
		resp, err := l.ClassStudents(req)
		httpx.FormatResponse(resp, err, w)
	}
}
