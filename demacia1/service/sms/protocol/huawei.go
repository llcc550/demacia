package protocol

import (
	"fmt"
	"strings"
)

var (
	HuaweiProtocolMsg = map[string]string{
		"000000":  "发送请求成功",
		"E000000": "系统异常, 请检查 templateParas 参数设置是否正确",
		"E000001": "HTTP 头中未找到 Authorization 字段",
		"E000002": "Authorization 字段中未找到 realm 字段",
		"E000003": "Authorization 字段中未找到 profile 字段",
		"E000004": "Authorization 字段中 realm 字段应该为 \"SDP\"",
		"E000005": "Authorization 字段中 profile 字段应该为 \"UsernameToken\"",
		"E000006": "Authorization 字段中 type 字段应该为 \"AppKey\"",
		"E000007": "Authorization 字段中 未找到 type 字段",
		"E000008": "Authorization 字段中 没有携带 WSSE",
		"E000020": "HTTP 头未找到 X-WSSE 字段",
		"E000021": "X-WSSE 字段中未找到 Username 字段",
		"E000022": "X-WSSE 字段中未找到 Nonce 字段",
		"E000023": "X-WSSE 字段中未找到 Created 字段",
		"E000024": "X-WSSE 字段中未找到 PasswordDigest 字段",
		"E000025": "Created 格式错误",
		"E000026": "X-WSSE 字段中未找到 UsernameToken属性",
		"E000027": "非法请求",
		"E000040": "ContentType 值应该为 application/www-form-urlencoded",
		"E000503": "参数格式错误",
		"E000510": "短信发送失败, 描述见参数 status",
		"E000623": "SP 短信发送量达到限额",
		"E000101": "鉴权失败",
		"E000102": "app_key 无效",
		"E000103": "app_key 不可用",
		"E000104": "app_secret 无效",
		"E000105": "PasswordDigest 无效",
		"E000106": "app_key 没有调用本 API 权限",
		"E000109": "用户状态未激活",
		"E000110": "时间超出限制",
		"E000111": "用户名或密码错误",
		"E000112": "用户状态已冻结",
		"E000620": "对端 app IP 不在白名单列表中",
	}
)

func NewHuaweiProtocol(data ApiProtocol) ApiProtocol {
	return data
}

type HuaweiProtocol struct {
	Result      []HuaweiResultProtocol `json:"result"`
	Code        string                 `json:"code"`
	Description string                 `json:"description"`
}

type HuaweiResultProtocol struct {
	OriginTo    string `json:"originTo"`
	CreatedTime string `json:"createdTime"`
	From        string `json:"from"`
	SmsMsgId    string `json:"smsMsgId"`
	Status      string `json:"status"`
}

func (p *HuaweiProtocol) GetCode() string {
	if p.Code == "000000" {
		return "0"
	}
	return fmt.Sprintf("%s", p.Code)
}

func (p *HuaweiProtocol) GetMsg() string {
	code := p.GetCode()
	if code == HuaweiApiServiceSuccess {
		return ""
	}
	var failedMsgList []string
	for _, prt := range p.Result {
		failedMsgList = append(failedMsgList, HuaweiProtocolMsg[prt.Status])
	}
	if len(failedMsgList) == 0 {
		return ""
	}
	return strings.Join(failedMsgList, ";")
}

func (p *HuaweiProtocol) GetData() interface{} {
	return p
}
