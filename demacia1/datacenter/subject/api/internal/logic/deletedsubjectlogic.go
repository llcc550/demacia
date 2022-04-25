package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/datacenter"
	"demacia/common/errlist"
	"demacia/datacenter/databus/rpc/databus"
	"demacia/datacenter/subject/api/internal/svc"
	"demacia/datacenter/subject/api/internal/types"
	"demacia/datacenter/subject/model"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type DeletedSubjectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletedSubjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeletedSubjectLogic {
	return DeletedSubjectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletedSubjectLogic) DeletedSubject(req types.Ids) error {

	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return errlist.NoAuth
	}
	if len(req.Id) > 0 {
		threading.GoSafe(func() {
			for _, v := range req.Id {
				err = l.svcCtx.SubjectModel.DeletedByOrgIdAndId(orgId, v)
				if err != nil {
					continue
				}
				_, _ = l.svcCtx.DataBusRpc.Delete(context.Background(), &databus.Req{
					Topic:    datacenter.Subject,
					ObjectId: v,
				})
				err = l.svcCtx.SubjectGradeModel.Deleted(&model.SubjectGrade{
					OrgId:     orgId,
					SubjectId: v,
				})
				if err != nil {
					continue
				}
			}
		})
	}
	return nil
}
