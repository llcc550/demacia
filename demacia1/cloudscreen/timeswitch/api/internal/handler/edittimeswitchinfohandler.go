package handler

import (
	"net/http"

	"demacia/cloudscreen/timeswitch/api/internal/logic"
	"demacia/cloudscreen/timeswitch/api/internal/svc"
	"demacia/cloudscreen/timeswitch/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func EditTimeSwitchInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TimeSwitchInfo
		if err := httpx.Parse(r, &req); err != nil {
			httpx.FormatResponse(nil, err, w)
			return
		}

		l := logic.NewEditTimeSwitchInfoLogic(r.Context(), ctx)
		resp, err := l.EditTimeSwitchInfo(req)
		httpx.FormatResponse(resp, err, w)
	}
}
