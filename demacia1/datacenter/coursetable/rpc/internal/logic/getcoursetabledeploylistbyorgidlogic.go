package logic

import (
	"context"
	"demacia/common/errlist"
	"demacia/datacenter/coursetable/errors"
	"demacia/datacenter/coursetable/rpc/coursetable"
	"demacia/datacenter/coursetable/rpc/internal/svc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GetCourseTableDeployListByOrgIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCourseTableDeployListByOrgIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCourseTableDeployListByOrgIdLogic {
	return &GetCourseTableDeployListByOrgIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCourseTableDeployListByOrgIdLogic) GetCourseTableDeployListByOrgId(in *coursetable.OrgIdReq) (*coursetable.CourseTableDeployResp, error) {

	if in.OrgId == 0 {
		return nil, status.Error(codes.InvalidArgument, errlist.InvalidParam.Error())
	}
	resp := &coursetable.CourseTableDeployResp{}
	deploys, err := l.svcCtx.CourseTableDeployModel.SelectByOrgId(in.OrgId)
	if err != nil {
		l.Logger.Errorf("select deployInfo err:%s", err.Error())
		return nil, status.Error(codes.NotFound, errors.DeployNotFoundErr.Error())
	}

	for _, deploy := range deploys {
		var courseFlag bool
		if deploy.CourseFlag == 1 {
			courseFlag = true
		}
		resp.List = append(resp.List, &coursetable.CourseTableDeploy{
			Note:       deploy.Note,
			Weekday:    int32(deploy.Weekday),
			StartTime:  deploy.StartTime,
			EndTime:    deploy.EndTime,
			Grouping:   int32(deploy.Grouping),
			CourseFlag: courseFlag,
			CourseSort: int32(deploy.CourseSort),
		})
	}

	return resp, nil
}
