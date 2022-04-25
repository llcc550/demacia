package logic

import (
	"context"
	"database/sql"

	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/class/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GetTeacherByClassIdAndSubjectIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTeacherByClassIdAndSubjectIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTeacherByClassIdAndSubjectIdLogic {
	return &GetTeacherByClassIdAndSubjectIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetTeacherByClassIdAndSubjectId 根据班级Id,学科Id获取任课关系
func (l *GetTeacherByClassIdAndSubjectIdLogic) GetTeacherByClassIdAndSubjectId(in *class.ClassSubjectIdReq) (*class.ClassSubjectTeachInfoResp, error) {
	var resp class.ClassSubjectTeachInfoResp
	findResp, err := l.svcCtx.ClassSubjectTeacherModel.FindTeacherByClassIdAndSubjectId(in.ClassId, in.SubjectId)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Error("Class Subject Teacher Rpc FindTeacherByClassIdAndSubjectId sql err:" + err.Error())
		return nil, err
	}
	resp = class.ClassSubjectTeachInfoResp{
		ClassId:      findResp.ClassId,
		SubjectId:    findResp.SubjectId,
		SubjectTitle: findResp.SubjectTitle,
		MemberId:     findResp.MemberId,
		TrueName:     findResp.TrueName,
	}
	return &resp, nil
}
