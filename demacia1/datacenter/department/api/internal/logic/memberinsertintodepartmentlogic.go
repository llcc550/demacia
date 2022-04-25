package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/department/api/internal/svc"
	"demacia/datacenter/department/api/internal/types"
	"demacia/datacenter/department/model"
	"demacia/datacenter/member/rpc/member"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/mr"
)

type MemberInsertIntoDepartmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMemberInsertIntoDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) MemberInsertIntoDepartmentLogic {
	return MemberInsertIntoDepartmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberInsertIntoDepartmentLogic) MemberInsertIntoDepartment(req types.MemberIdsReq) error {
	if req.DepartmentId <= 0 || len(req.MemberIds) == 0 {
		return errlist.InvalidParam
	}
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return errlist.AuthLoginFail
	}
	departmentInfo, err := l.svcCtx.DepartmentModel.GetDepartmentById(req.DepartmentId)
	if err != nil || departmentInfo.OrgId != orgId {
		return errlist.DepartmentNotExit
	}
	list, err := l.svcCtx.DepartmentMemberModel.GetMembersByOrgIdAndDepartmentId(orgId, req.DepartmentId)
	if err != nil {
		return err
	}
	memberIdMap := map[int64]interface{}{}
	for _, memberInfo := range list {
		memberIdMap[memberInfo.MemberId] = nil
	}
	memberIds := make([]int64, 0, len(req.MemberIds))
	i := int64(0)
	for _, memberId := range req.MemberIds {
		if _, ok := memberIdMap[memberId]; !ok {
			memberIds = append(memberIds, memberId)
		} else {
			i++
		}
	}
	members := make(model.DepartmentMembers, 0, len(memberIds))
	err = mr.MapReduceVoid(func(source chan<- interface{}) {
		for _, memberId := range memberIds {
			source <- memberId
		}
	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
		memberInfo, err := l.svcCtx.MemberRpc.FindOneById(l.ctx, &member.IdReq{Id: item.(int64)})
		if err != nil || memberInfo.OrgId != orgId {
			return
		}
		writer.Write(&model.DepartmentMember{
			MemberId:     memberInfo.Id,
			DepartmentId: req.DepartmentId,
			OrgId:        orgId,
			TrueName:     memberInfo.TrueName,
			Mobile:       memberInfo.Mobile,
		})
	}, func(pipe <-chan interface{}, cancel func(error)) {
		for item := range pipe {
			members = append(members, item.(*model.DepartmentMember))
		}
	})
	if err != nil {
		return err
	}
	err = l.svcCtx.DepartmentMemberModel.BatchInsert(members)
	if err != nil {
		return err
	}
	err = l.svcCtx.DepartmentModel.UpdateDepartmentMemberCount(req.DepartmentId, int64(len(members))+i)
	if err != nil {
		return err
	}
	return nil
}
