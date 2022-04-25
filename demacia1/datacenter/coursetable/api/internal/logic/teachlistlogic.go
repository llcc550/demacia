package logic

import (
	"context"
	"database/sql"
	"demacia/common/errlist"
	"demacia/datacenter/class/rpc/class"

	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"

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
	clazz, err := l.svcCtx.ClassRpc.GetClassInfoById(l.ctx, &class.IdReq{Id: req.ClassId})
	if err != nil && err == sql.ErrNoRows {
		l.Logger.Errorf("Class Api TeachList class[model] GetClassById err：%s", err.Error())
		return nil, errlist.Unknown
	}
	if clazz == nil {
		return nil, errlist.ClassNotFound
	}
	list, err := l.svcCtx.TeachModel.ListByClassId(clazz.Id)
	if err != nil {
		l.Logger.Errorf("Class Api TeachList Teach[model] ListByClassId err：%s", err.Error())
		return nil, errlist.Unknown
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
