package logic

import (
	"context"
	"encoding/json"

	"demacia/common/datacenter"
	"demacia/datacenter/department/model"
	"demacia/datacenter/department/rmq/internal/svc"
	"demacia/datacenter/member/rpc/member"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type Consumer struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConsumerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Consumer {
	return &Consumer{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *Consumer) MemberConsume(_, v string) {
	var message datacenter.Message
	err := json.Unmarshal([]byte(v), &message)
	if err != nil {
		return
	}
	memberId := message.ObjectId
	memberInfo, err := l.svcCtx.MemberRpc.FindOneById(l.ctx, &member.IdReq{Id: memberId})
	if err != nil {
		return
	}
	orgId := memberInfo.OrgId

	switch message.Action {
	case datacenter.Delete:
		list, err := l.svcCtx.DepartmentMemberModel.GetDepartmentIdsByOrgIdAndMemberId(orgId, memberId)
		if err != nil || len(list) == 0 {
			return
		}
		_ = l.svcCtx.DepartmentMemberModel.DeleteByByOrgIdAndMemberId(orgId, memberId)
		for _, item := range list {
			members, err := l.svcCtx.DepartmentMemberModel.GetMembersByOrgIdAndDepartmentId(orgId, item)
			if err != nil {
				continue
			}
			_ = l.svcCtx.DepartmentModel.UpdateDepartmentMemberCount(item, int64(len(members)))
		}
	case datacenter.Update:
		_ = l.svcCtx.DepartmentMemberModel.UpdateMemberInfo(&model.DepartmentMember{
			MemberId: memberId,
			OrgId:    orgId,
			TrueName: memberInfo.TrueName,
			Mobile:   memberInfo.Mobile,
		})
	}
}
