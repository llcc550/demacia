package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/datacenter/card/rpc/card"
	"demacia/datacenter/student/api/internal/svc"
	"demacia/datacenter/student/api/internal/types"
	"demacia/datacenter/student/errors"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type StudentsDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStudentsDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) StudentsDetailLogic {
	return StudentsDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StudentsDetailLogic) StudentsDetail(req types.IdRequest) (resp *types.StudentDetail, err error) {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, err
	}
	studentInfo, err := l.svcCtx.StudentModel.FindOneById(req.StudentId)
	if err != nil || orgId != studentInfo.OrgId {
		return resp, errors.StudentNotExist
	}
	cardSlice := make([]string, 0)
	cardNumber, err := l.svcCtx.CardRpc.GetStudentCardList(l.ctx, &card.ListReq{
		OrgId:    orgId,
		ObjectId: req.StudentId,
	})
	if err != nil {
		logx.Errorf("card err:%s", err.Error())
	} else {
		cardSlice = cardNumber.CardNum
	}

	resp = &types.StudentDetail{
		StudentId:   studentInfo.Id,
		ClassId:     studentInfo.ClassId,
		ClassName:   studentInfo.ClassFullName,
		UserName:    studentInfo.UserName,
		StudentName: studentInfo.TrueName,
		Sex:         studentInfo.Sex,
		CardNumber:  cardSlice,
		IdNumber:    studentInfo.IdNumber,
		Avatar:      studentInfo.Avatar,
		Face:        studentInfo.Face,
	}
	return resp, nil
}
