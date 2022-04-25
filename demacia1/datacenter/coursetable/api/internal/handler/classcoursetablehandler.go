package handler

import (
	"demacia/common/errlist"
	"demacia/datacenter/coursetable/api/internal/logic"
	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"
	"fmt"
	"net/http"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func ClassCourseTableHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ClassIdReq
		if err := httpx.Parse(r, &req); err != nil {
			fmt.Println(err)
			httpx.FormatResponse(nil, errlist.InvalidParam, w)
			return
		}

		l := logic.NewClassCourseTableLogic(r.Context(), ctx)
		resp, err := l.ClassCourseTable(req)
		httpx.FormatResponse(resp, err, w)
	}
}
