package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"demacia/datacenter/student/api/internal/logic"
	"demacia/datacenter/student/api/internal/svc"
	"demacia/datacenter/student/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

func exportStudentsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListConditionRequest
		if err := httpx.Parse(r, &req); err != nil {
			fmt.Println(err)
			httpx.FormatResponse(nil, err, w)
			return
		}
		l := logic.NewExportStudentsLogic(r.Context(), ctx)
		resp, _ := l.ExportStudents(req)
		w.Header().Add("Content-Type", "applicationnd.ms-excel")
		w.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
		w.Header().Add("Content-Disposition", "attachment;filename="+url.QueryEscape(resp.Name)+".xlsx")
		w.Header().Add("Cache-Control", "max-age=0")
		_ = resp.File.Write(w)
	}
}
