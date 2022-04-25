package handler

import (
	"demacia/cloudscreen/courserecord/api/internal/logic"
	"demacia/cloudscreen/courserecord/api/internal/svc"
	"demacia/cloudscreen/courserecord/api/internal/types"
	"net/http"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func courseRecordCountHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseRecordCountReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewCourseRecordCountLogic(r.Context(), ctx)
		resp, err := l.CourseRecordCount(req)
		httpx.FormatResponse(resp, err, w)
	}
}
