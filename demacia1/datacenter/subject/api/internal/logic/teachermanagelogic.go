package logic

import (
	"context"
	"database/sql"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/member/rpc/member"
	"demacia/datacenter/subject/api/internal/svc"
	"demacia/datacenter/subject/api/internal/types"
	"demacia/datacenter/subject/model"
	"fmt"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type TeacherManageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTeacherManageLogic(ctx context.Context, svcCtx *svc.ServiceContext) TeacherManageLogic {
	return TeacherManageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TeacherManageLogic) TeacherManage(req types.TeacherManageReq) error {

	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return errlist.NoAuth
	}
	subject, err := l.svcCtx.SubjectModel.FindSubjectById(req.Id)
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Subject Api TeacherManage subject[model] FindSubjectById err : %s ", err.Error())
		return errlist.Unknown
	}
	if subject == nil {
		return errlist.SubjectNotFound
	}
	threading.GoSafe(func() {
		members := make([]types.SubjectTeacherResp, 0)
		for _, teacher := range req.Teachers {
			memberInfo, err := l.svcCtx.MemberRpc.FindOneById(context.Background(), &member.IdReq{Id: teacher})
			if err != nil {
				l.Logger.Errorf("Subject Api Member[Rpc] FindOneById err :%s", err.Error())
				continue
			}
			members = append(members, types.SubjectTeacherResp{
				Id:       memberInfo.Id,
				TrueName: memberInfo.TrueName,
			})
		}

		subjectTeacher := make(model.SubjectTeachers, 0)

		if len(members) > 0 {
			for _, v := range members {
				fmt.Println(v)
				subjectTeacher = append(subjectTeacher, &model.SubjectTeacher{
					SubjectId: subject.Id,
					MemberId:  v.Id,
					TrueName:  v.TrueName,
					OrgId:     orgId,
				})
			}
		}
		if len(subjectTeacher) > 0 {
			err := l.svcCtx.SubjectTeacherModel.DeletedBySubjectId(subject.Id)
			if err != nil {
				l.Logger.Errorf("TeacherManage Deleted Subject err:%s", err.Error())
			}
			err = l.svcCtx.SubjectTeacherModel.BatchInsert(subjectTeacher)
			if err != nil {
				l.Logger.Errorf("TeacherManage BatchInsert Subject err:%s", err.Error())
			}
		}
	})
	return nil
}
