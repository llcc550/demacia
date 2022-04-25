package handler

import (
	"demacia/common/errlist"
	"net/http"

	"demacia/datacenter/coursetable/api/internal/logic"
	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func CourseTableAddHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseTableAddReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, errlist.InvalidParam, w)
			return
		}

		l := logic.NewCourseTableAddLogic(r.Context(), ctx)
		resp, err := l.CourseTableAdd(req)
		httpx.FormatResponse(resp, err, w)
	}
}
