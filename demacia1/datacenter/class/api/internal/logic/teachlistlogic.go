package logic

import (
	"context"
	"demacia/common/errlist"
	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type TeachListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTeachListLogic(ctx context.Context, svcCtx *svc.ServiceContext) TeachListLogic {
	return TeachListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TeachListLogic) TeachList(req types.Id) (resp *types.ListTeachResp, err error) {

	if req.Id < 0 {
		return nil, errlist.ClassNotFound
	}
	class, err := l.svcCtx.ClassModel.GetClassById(req.Id)
	if err != nil || class == nil {
		return nil, errlist.ClassNotFound
	}
	list, err := l.svcCtx.TeachModel.ListByClassId(class.Id)
	if err != nil {
		return nil, err
	}
	resp = &types.ListTeachResp{List: []types.TeachResp{}}
	for _, v := range list {
		resp.List = append(resp.List, types.TeachResp{
			Id:           v.Id,
			ClassId:      v.ClassId,
			SubjectId:    v.SubjectId,
			SubjectTitle: v.SubjectTitle,
			TeacherId:    v.MemberId,
			TeacherName:  v.TrueName,
		})
	}
	return resp, nil
}
