package handler

import (
	"demacia/datacenter/subject/api/internal/logic"
	"demacia/datacenter/subject/api/internal/svc"
	"demacia/datacenter/subject/api/internal/types"
	"net/http"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func GetSubjectTeacherHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Id
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewGetSubjectTeacherLogic(r.Context(), ctx)
		resp, err := l.GetSubjectTeacher(req)
		httpx.FormatResponse(resp, err, w)
	}
}
