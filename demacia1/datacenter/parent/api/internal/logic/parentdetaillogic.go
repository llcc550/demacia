package logic

import (
	"context"

	"demacia/datacenter/parent/api/internal/svc"
	"demacia/datacenter/parent/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type ParentDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewParentDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) ParentDetailLogic {
	return ParentDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ParentDetailLogic) ParentDetail(req types.IdRequest) (resp *types.DetailResponse, err error) {
	studentParentInfo, err := l.svcCtx.StudentParentModel.FindListByParentId(req.ParentId)
	if err != nil {
		return nil, err
	}
	studentParent := make([]types.StudentInfo, 0, len(*studentParentInfo))
	for _, info := range *studentParentInfo {
		studentParent = append(studentParent, types.StudentInfo{
			StudentId:   info.StudentId,
			StudentName: info.StudentName,
			ClassName:   info.ClassName,
			Relation:    info.Relation,
		})
	}
	resp = &types.DetailResponse{
		ParentId:    0,
		StudentInfo: []types.StudentInfo{},
		ParentName:  "",
		Moblie:      "",
		IdNumber:    "",
		Address:     "",
		Face:        "",
	}
	parentInfo, err := l.svcCtx.ParentModel.FindOneById(req.ParentId)
	if err != nil {
		return nil, err
	}
	resp.ParentId = parentInfo.Id
	resp.StudentInfo = studentParent
	resp.ParentName = parentInfo.TrueName
	resp.Face = parentInfo.Face
	resp.IdNumber = parentInfo.IdNumber
	resp.Moblie = parentInfo.Mobile
	resp.Address = parentInfo.Address
	return resp, nil
}
