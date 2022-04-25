package handler

import (
	"net/http"

	"demacia/service/websocket/api/internal/logic"
	"demacia/service/websocket/api/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func nodesHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewNodesLogic(r.Context(), ctx)
		resp, err := l.Nodes()
		httpx.FormatResponse(resp, err, w)
	}
}
