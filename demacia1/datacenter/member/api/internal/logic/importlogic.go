package logic

import (
	"context"
	"fmt"
	"os"
	"strings"

	"demacia/common/baseauth"
	"demacia/common/basefunc"
	"demacia/common/datacenter"
	"demacia/common/errlist"
	"demacia/datacenter/databus/rpc/databus"
	"demacia/datacenter/member/api/internal/svc"
	"demacia/datacenter/member/api/internal/types"
	"demacia/datacenter/member/model"
	"demacia/datacenter/organization/rpc/organization"
	"demacia/service/websocket/rpc/websocket"

	"github.com/tealeg/xlsx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

type ImportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) ImportLogic {
	return ImportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImportLogic) Import(req types.ImportReq) error {
	OrgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	// 验证机构是否存在
	org, err := l.svcCtx.OrganizationRpc.FindOne(l.ctx, &organization.IdReply{Id: OrgId})
	if err != nil {
		return err
	}

	threading.GoSafe(func() {
		l.importMembers(org.Id, req.Url, req.WebsocketUuid)
	})
	return nil
}

func (l *ImportLogic) importMembers(Oid int64, url string, uuid string) {
	fn := func(err *httpx.CodeError, msg ...interface{}) {
		_, _ = l.svcCtx.WebsocketRpc.Push(context.Background(), &websocket.Request{
			Key:  uuid,
			Code: int64(err.Code()),
			Msg:  fmt.Sprintf(err.Error(), msg...),
		})
	}
	path := basefunc.Download(url)
	if path == "" {
		fn(errlist.Excel)
		return
	}
	defer os.Remove(path)
	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		fn(errlist.Excel)
		return
	}

	for _, sheet := range xlFile.Sheets { // 遍历sheet
		for k, r := range sheet.Rows { // 遍历每一行
			if k == 0 { // 过滤掉第一行
				continue
			}
			// Cells 遍历每一行中的每一个单元格
			if len(r.Cells) < 3 {
				break
			}

			trueName := strings.TrimSpace(r.Cells[0].Value) // 昵称
			mobile := strings.TrimSpace(r.Cells[1].Value)   // 手机号
			sexText := strings.TrimSpace(r.Cells[2].Value)  // 性别

			// 如果存在字段为空的就跳出
			if trueName == "" || mobile == "" || sexText == "" {
				break
			}
			// 查询机构是否存在

			info := model.Member{
				OrgId:    Oid,
				TrueName: strings.TrimSpace(r.Cells[0].Value),
				UserName: strings.TrimSpace(r.Cells[1].Value),
				Password: basefunc.HashPassword(strings.TrimSpace(r.Cells[1].Value)),
				Mobile:   strings.TrimSpace(r.Cells[1].Value),
			}
			// 手机号格式
			if !basefunc.CheckMobile(info.Mobile) {
				fn(errlist.ImportErr, info.TrueName, "用户手机号格式错误")
				continue
			}
			// 手机号唯一
			byMobile, _ := l.svcCtx.MemberModel.FindOneByMobile(info.Mobile, Oid)
			if byMobile != nil {
				fn(errlist.ImportErr, info.TrueName, "用户手机号已存")
				continue
			}

			switch sexText {
			case "男":
				info.Sex = 1
			case "女":
				info.Sex = 0
			default:
				info.Sex = 1
			}
			memberId, err := l.svcCtx.MemberModel.Insert(&info)
			if err != nil {
				fn(errlist.ImportErr, info.TrueName, "新增失败")
				continue
			}
			_, _ = l.svcCtx.DataBusRpc.Create(context.Background(), &databus.Req{Topic: datacenter.Member, ObjectId: memberId})
		}
	}
}
