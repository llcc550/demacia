package handler

import (
	"demacia/datacenter/coursetable/api/internal/logic"
	"demacia/datacenter/coursetable/api/internal/svc"
	"net/http"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func OrgCourseTableHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewOrgCourseTableLogic(r.Context(), ctx)
		resp, err := l.OrgCourseTable()
		httpx.FormatResponse(resp, err, w)
	}
}
