package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"
	"demacia/datacenter/coursetable/errors"
	"demacia/datacenter/coursetable/model"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"time"
)

type CourseTableDeploySaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseTableDeploySaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) CourseTableDeploySaveLogic {
	return CourseTableDeploySaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseTableDeploySaveLogic) CourseTableDeploySave(req types.CourseTableDeploySaveReq) (*types.BoolReply, error) {

	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return &types.BoolReply{Success: false}, errlist.NoAuth
	}

	if len(req.DeployInfo) == 0 && len(req.DeployInfo) > 7 {
		return &types.BoolReply{Success: false}, errlist.InvalidParam
	}

	for i, info := range req.DeployInfo {
		for j := i; j < len(req.DeployInfo); j++ {
			if info.Weekday == req.DeployInfo[j].Weekday {
				startTime, err := time.Parse("15:04:05", info.StartTime)
				if err != nil || startTime.IsZero() {
					return &types.BoolReply{Success: false}, errors.InvalidTime
				}
				endTime, err := time.Parse("15:04:05", req.DeployInfo[j].EndTime)
				if err != nil || endTime.IsZero() {
					return &types.BoolReply{Success: false}, errors.InvalidTime
				}
				if startTime.After(endTime) {
					return &types.BoolReply{Success: false}, errors.CourseTimeErr
				}
			}
		}
		for j := i + 1; j < len(req.DeployInfo); j++ {
			if info.Weekday == req.DeployInfo[j].Weekday {
				endTime, err := time.Parse("15:04:05", info.EndTime)
				if err != nil || endTime.IsZero() {
					return &types.BoolReply{Success: false}, errors.InvalidTime
				}
				startTime, err := time.Parse("15:04:05", req.DeployInfo[j].StartTime)
				if err != nil || startTime.IsZero() {
					return &types.BoolReply{Success: false}, errors.InvalidTime
				}
				if endTime.After(startTime) {
					return &types.BoolReply{Success: false}, errors.CourseTimeErr
				}
			}
		}
	}

	var deploys model.CourseTableDeploys

	for _, deploy := range req.DeployInfo {
		deploys = append(deploys, &model.CourseTableDeploy{
			OrgId:      oid,
			CourseSort: deploy.CourseSort,
			Weekday:    deploy.Weekday,
			Note:       deploy.Note,
			Grouping:   deploy.Grouping,
			StartTime:  deploy.StartTime,
			EndTime:    deploy.EndTime,
			CourseFlag: deploy.CourseFlag,
		})
	}

	err = l.svcCtx.CourseTableDeployModel.DeleteByOrgId(oid)
	if err != nil {
		return &types.BoolReply{Success: false}, errlist.Unknown
	}

	err = l.svcCtx.CourseTableDeployModel.InsertCourseTableDeploy(deploys)
	if err != nil {
		return &types.BoolReply{Success: false}, errlist.Unknown
	}

	return &types.BoolReply{Success: true}, nil
}
