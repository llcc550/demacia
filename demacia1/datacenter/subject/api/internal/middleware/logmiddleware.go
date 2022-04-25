package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"demacia/common/baseauth"
	"demacia/datacenter/databus/rpc/databusclient"

	"github.com/thinkeridea/go-extend/exnet"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type LogMiddleware struct {
	DataBusRpc databusclient.Databus
}

func NewLogMiddleware(dataBus databusclient.Databus) *LogMiddleware {
	return &LogMiddleware{
		DataBusRpc: dataBus,
	}
}

func (m *LogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := exnet.ClientPublicIP(r)
		jwt, _ := baseauth.GetUserJwt(r)
		jwtJson, _ := json.Marshal(jwt)
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(r.Body)
		threading.GoSafe(func() {
			_, _ = m.DataBusRpc.Log(context.Background(), &databusclient.LogReq{Ip: ip, Route: r.RequestURI, Jwt: string(jwtJson), Req: buf.String()})
		})
		r.Body = ioutil.NopCloser(buf)
		next(w, r)
	}
}
