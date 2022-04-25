package handler

import (
	"demacia/datacenter/subject/api/internal/logic"
	"demacia/datacenter/subject/api/internal/svc"
	"demacia/datacenter/subject/api/internal/types"
	"fmt"
	"net/http"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func AddSubjectHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddSubjectReq
		if err := httpx.Parse(r, &req); err != nil {
			fmt.Println(err)
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewAddSubjectLogic(r.Context(), ctx)
		resp, err := l.AddSubject(req)
		httpx.FormatResponse(resp, err, w)
	}
}
