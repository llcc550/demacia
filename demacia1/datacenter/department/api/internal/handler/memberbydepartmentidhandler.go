package handler

import (
	"net/http"

	"demacia/datacenter/department/api/internal/logic"
	"demacia/datacenter/department/api/internal/svc"
	"demacia/datacenter/department/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func memberByDepartmentIdHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MemberListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewMemberByDepartmentIdLogic(r.Context(), ctx)
		resp, err := l.MemberByDepartmentId(req)
		httpx.FormatResponse(resp, err, w)
	}
}
