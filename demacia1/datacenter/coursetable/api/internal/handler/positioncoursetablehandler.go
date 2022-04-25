package handler

import (
	"demacia/common/errlist"
	"demacia/datacenter/coursetable/api/internal/logic"
	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"
	"net/http"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func PositionCourseTableHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PositionIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, errlist.InvalidParam, w)
			return
		}

		l := logic.NewPositionCourseTableLogic(r.Context(), ctx)
		resp, err := l.PositionCourseTable(req)
		httpx.FormatResponse(resp, err, w)
	}
}
