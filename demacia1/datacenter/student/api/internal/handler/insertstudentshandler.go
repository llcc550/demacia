package handler

import (
	"net/http"

	"demacia/datacenter/student/api/internal/logic"
	"demacia/datacenter/student/api/internal/svc"
	"demacia/datacenter/student/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func insertStudentsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.InsertRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}
		l := logic.NewInsertStudentsLogic(r.Context(), ctx)
		err := l.InsertStudents(req)
		httpx.FormatResponse(nil, err, w)
	}
}
