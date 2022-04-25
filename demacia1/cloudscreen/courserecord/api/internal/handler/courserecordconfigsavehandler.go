package handler

import (
	"demacia/cloudscreen/courserecord/api/internal/logic"
	"demacia/cloudscreen/courserecord/api/internal/svc"
	"demacia/cloudscreen/courserecord/api/internal/types"
	"demacia/common/errlist"
	"net/http"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func courseRecordConfigSaveHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseRecordConfigSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, errlist.InvalidParam, w)
			return
		}

		l := logic.NewCourseRecordConfigSaveLogic(r.Context(), ctx)
		resp, err := l.CourseRecordConfigSave(req)
		httpx.FormatResponse(resp, err, w)
	}
}
