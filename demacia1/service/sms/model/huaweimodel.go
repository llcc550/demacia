package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"demacia/service/sms/protocol"

	"github.com/google/uuid"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type (
	HuaweiConfig struct {
		BaseSendUrl  string
		BatchSendUrl string
		AppKey       string
		AppSecret    string
		Templates    map[string]Template
	}

	Template struct {
		TemplateSign string
		TemplateId   string
	}

	HuaweiModel struct {
		config *HuaweiConfig
		limit  uint8
	}

	ContentMobile struct {
		TemplateId string   `json:"template_id"`
		Mobile     []string `json:"mobile"`
		Params     []string `json:"params"`
	}

	HuaweiRequest struct {
		Sign          string
		To            string
		Types         uint
		TemplateId    string
		TemplateParas []string
		MultiMt       []*ContentMobile
	}

	multiMt struct {
		From       string       `json:"from"`
		SmsContent []smsContent `json:"smsContent"`
	}

	smsContent struct {
		To            []string `json:"to"`
		TemplateId    string   `json:"templateId"`
		TemplateParas []string `json:"templateParas"`
	}

	HuaweiResponse struct {
		protocol.HuaweiProtocol
	}
)

var (
	ErrMobileNumberOutLimit = errors.New("手机号码个数超过1000条限制")
	HuaweiSmsAcceptCodes    = []int{
		http.StatusOK,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
	}
)

func (p *HuaweiResponse) GetData() interface{} {
	return p
}

func (hwr *HuaweiModel) Send(mobile, templateId string, params []string) error {
	var resp HuaweiResponse
	var req = &HuaweiRequest{
		Sign:          hwr.config.Templates[templateId].TemplateSign,
		To:            mobile,
		TemplateParas: params,
		TemplateId:    hwr.config.Templates[templateId].TemplateId,
	}

	err := hwr.SendSmsApi(req, &resp)
	if err != nil {
		logx.Error(err)
	}

	if resp.Code != "000000" {
		logx.Errorf("send sms failed, error: %s", resp.GetMsg())
	}
	return err
}

func (hwr *HuaweiModel) GetTimeUTC() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

func NewHuaweiModel(config *HuaweiConfig) *HuaweiModel {
	return &HuaweiModel{
		config: config,
	}
}

func (hwr *HuaweiModel) MultiSendBatchSms(templateId string, req []*ContentMobile) error {
	if len(req) > 1000 {
		return ErrMobileNumberOutLimit
	}
	var resp HuaweiResponse
	var request = &HuaweiRequest{
		Types:   2,
		Sign:    hwr.config.Templates[templateId].TemplateSign,
		MultiMt: req,
	}
	return hwr.SendSmsApi(request, &resp)
}

func (hwr *HuaweiModel) SendSmsApi(req *HuaweiRequest, resp *HuaweiResponse) error {
	timeUTC := hwr.GetTimeUTC()
	nonce := uuid.New().String()
	sha256Hash := sha256.New()
	sha256Hash.Write([]byte(nonce + timeUTC + hwr.config.AppSecret))
	wsse := fmt.Sprintf("UsernameToken Username=\"%s\",PasswordDigest=\"%s\",Nonce=\"%s\",Created=\"%s\"",
		hwr.config.AppKey,
		base64.StdEncoding.EncodeToString(sha256Hash.Sum(nil)),
		nonce,
		timeUTC,
	)
	body := url.Values{
		"from":          {req.Sign},
		"to":            {req.To},
		"templateId":    {req.TemplateId},
		"templateParas": {"[" + strings.Join(req.TemplateParas, ",") + "]"},
	}
	protocolClient := &protocol.ApiOption{
		Url:    hwr.config.BaseSendUrl,
		Method: http.MethodPost,
		Data:   protocol.NewHuaweiProtocol(resp),
		Header: map[string]string{
			protocol.HeaderContentTypeKey: protocol.ContentTypeForm,
			protocol.HeaderAcceptKey:      protocol.AcceptJson,
			protocol.HeaderWSSE:           wsse,
			protocol.HeaderAuthorization:  "WSSE  realm=\"SDP\",profile=\"UsernameToken\",type=\"Appkey\"",
		},
		Body:               bytes.NewBuffer([]byte(body.Encode())),
		AcceptCodes:        HuaweiSmsAcceptCodes,
		InsecureSkipVerify: true,
	}
	if req.Types > 0 {
		protocolClient.Url = hwr.config.BatchSendUrl
		protocolClient.Header[protocol.HeaderContentTypeKey] = protocol.ContentTypeJson
		var smsContents []smsContent
		for _, mobile := range req.MultiMt {
			smsContents = append(smsContents, smsContent{
				To:            mobile.Mobile,
				TemplateId:    hwr.config.Templates[mobile.TemplateId].TemplateId,
				TemplateParas: mobile.Params,
			})
		}
		jsonBody, _ := json.Marshal(multiMt{
			From:       req.Sign,
			SmsContent: smsContents,
		})
		protocolClient.Body = bytes.NewBuffer(jsonBody)
	}

	err := (protocolClient).Get()
	logx.Infof("sms protocol resp is : %+v", resp)
	return err
}
