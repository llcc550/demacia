package handler

import (
	"demacia/cloudscreen/courserecord/api/internal/logic"
	"demacia/cloudscreen/courserecord/api/internal/svc"
	"demacia/cloudscreen/courserecord/api/internal/types"
	"demacia/common/errlist"
	"net/http"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func GetCourseRecordListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseRecordReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, errlist.InvalidParam, w)
			return
		}

		l := logic.NewGetCourseRecordListLogic(r.Context(), ctx)
		resp, err := l.GetCourseRecordList(req)
		httpx.FormatResponse(resp, err, w)
	}
}
