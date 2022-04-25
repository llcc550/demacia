package handler

import (
	"demacia/cloudscreen/courserecord/api/internal/logic"
	"demacia/cloudscreen/courserecord/api/internal/svc"
	"net/http"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func courseRecordConfigInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewCourseRecordConfigInfoLogic(r.Context(), ctx)
		resp, err := l.CourseRecordConfigInfo()
		httpx.FormatResponse(resp, err, w)
	}
}
