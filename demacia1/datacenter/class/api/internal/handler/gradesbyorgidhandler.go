package handler

import (
	"net/http"

	"demacia/datacenter/class/api/internal/logic"
	"demacia/datacenter/class/api/internal/svc"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func gradesByOrgIdHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGradesByOrgIdLogic(r.Context(), ctx)
		resp, err := l.GradesByOrgId()
		httpx.FormatResponse(resp, err, w)
	}
}
