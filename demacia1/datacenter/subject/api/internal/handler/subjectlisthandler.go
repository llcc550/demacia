package handler

import (
	"net/http"

	"demacia/datacenter/subject/api/internal/logic"
	"demacia/datacenter/subject/api/internal/svc"
	"demacia/datacenter/subject/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func SubjectListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TitleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewSubjectListLogic(r.Context(), ctx)
		resp, err := l.SubjectList(req)
		httpx.FormatResponse(resp, err, w)
	}
}
