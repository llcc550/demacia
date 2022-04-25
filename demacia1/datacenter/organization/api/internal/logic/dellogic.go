package logic

import (
	"context"

	"demacia/common/datacenter"
	"demacia/common/errlist"
	"demacia/datacenter/databus/rpc/databus"
	"demacia/datacenter/member/rpc/memberclient"
	"demacia/datacenter/organization/api/internal/errors"
	"demacia/datacenter/organization/api/internal/svc"
	"demacia/datacenter/organization/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/mr"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type DelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelLogic(ctx context.Context, svcCtx *svc.ServiceContext) DelLogic {
	return DelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelLogic) Del(req types.IdReq) error {
	info, err := l.svcCtx.OrganizationModel.FindOneById(req.Id)
	if err != nil {
		return errors.OrganizationNotExist
	}
	err = mr.Finish(func() error {
		err := l.svcCtx.OrganizationModel.DeleteById(req.Id)
		if err != nil {
			l.Logger.Errorf("delete organization error. id is %d, error is %s", req.Id, err.Error())
		}
		return err
	}, func() error {
		_, err := l.svcCtx.MemberRpc.DeleteById(l.ctx, &memberclient.IdReq{Id: info.ManagerMemberId})
		if err != nil {
			l.Logger.Errorf("delete organization error when delete manager member. org_id is %d,member_id is %d, error is %s", req.Id, info.ManagerMemberId, err.Error())
		}
		return err
	})
	if err != nil {
		return errlist.Unknown
	}
	threading.GoSafe(func() {
		_, _ = l.svcCtx.DataBusRpc.Delete(context.Background(), &databus.Req{Topic: datacenter.Organization, ObjectId: req.Id})
	})
	return nil
}
