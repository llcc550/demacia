package handler

import (
	"demacia/cloudscreen/courserecord/api/internal/logic"
	"demacia/cloudscreen/courserecord/api/internal/svc"
	"demacia/cloudscreen/courserecord/api/internal/types"
	"net/http"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func courseRecordAddHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseRecordAddReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewCourseRecordAddLogic(r.Context(), ctx)
		resp, err := l.CourseRecordAdd(req)
		httpx.FormatResponse(resp, err, w)
	}
}
