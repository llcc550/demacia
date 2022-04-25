package handler

import (
	"demacia/common/errlist"
	"fmt"
	"net/http"

	"demacia/datacenter/coursetable/api/internal/logic"
	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func CourseTableDeploySaveHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseTableDeploySaveReq
		if err := httpx.Parse(r, &req); err != nil {
			fmt.Println(err)
			httpx.FormatResponse(nil, errlist.InvalidParam, w)
			return
		}

		l := logic.NewCourseTableDeploySaveLogic(r.Context(), ctx)
		resp, err := l.CourseTableDeploySave(req)
		httpx.FormatResponse(resp, err, w)
	}
}
