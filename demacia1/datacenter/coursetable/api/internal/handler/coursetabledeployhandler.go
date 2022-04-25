package handler

import (
	"net/http"

	"demacia/datacenter/coursetable/api/internal/logic"
	"demacia/datacenter/coursetable/api/internal/svc"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func CourseTableDeployHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewCourseTableDeployLogic(r.Context(), ctx)
		resp, err := l.CourseTableDeploy()
		httpx.FormatResponse(resp, err, w)
	}
}
