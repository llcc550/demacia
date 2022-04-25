package handler

import (
	"demacia/datacenter/coursetable/api/internal/logic"
	"demacia/datacenter/coursetable/api/internal/svc"
	"net/http"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func ClearCourseDeployHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewClearCourseDeployLogic(r.Context(), ctx)
		resp, err := l.ClearCourseDeploy()
		httpx.FormatResponse(resp, err, w)
	}
}
