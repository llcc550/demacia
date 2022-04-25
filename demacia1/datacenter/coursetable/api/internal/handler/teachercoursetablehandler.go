package handler

import (
	"demacia/common/errlist"
	"demacia/datacenter/coursetable/api/internal/logic"
	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"
	"net/http"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func TeacherCourseTableHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MemberIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, errlist.InvalidParam, w)
			return
		}

		l := logic.NewTeacherCourseTableLogic(r.Context(), ctx)
		resp, err := l.TeacherCourseTable(req)
		httpx.FormatResponse(resp, err, w)
	}
}
