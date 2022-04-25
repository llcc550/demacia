package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"time"

	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type CourseTableDeployLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseTableDeployLogic(ctx context.Context, svcCtx *svc.ServiceContext) CourseTableDeployLogic {
	return CourseTableDeployLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseTableDeployLogic) CourseTableDeploy() (*types.CourseTableDeployReply, error) {

	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return &types.CourseTableDeployReply{DeployInfo: []*types.DeployInfo{}}, errlist.InvalidParam
	}

	courseTableDeploys, err := l.svcCtx.CourseTableDeployModel.SelectByOrgId(oid)
	if err != nil || len(courseTableDeploys) == 0 {
		courseTableDeploys, err = l.svcCtx.CourseTableDeployModel.SelectByOrgId(0)
		if err != nil {
			logx.Errorf("select course_table_deploy err:%s", err)
			return &types.CourseTableDeployReply{DeployInfo: []*types.DeployInfo{}}, errlist.Unknown
		}
	}

	var res types.CourseTableDeployReply

	for i := 1; i < 8; i++ {
		res.DeployInfo = append(res.DeployInfo, &types.DeployInfo{
			Weekday:          int8(i),
			CourseDeployInfo: []*types.CourseDeployInfo{},
		})
	}

	for _, deploy := range courseTableDeploys {
		for _, info := range res.DeployInfo {
			if deploy.Weekday == info.Weekday {
				startTime, _ := time.Parse("2006-01-02T15:04:05Z", deploy.StartTime)
				endTime, _ := time.Parse("2006-01-02T15:04:05Z", deploy.EndTime)
				info.CourseDeployInfo = append(info.CourseDeployInfo, &types.CourseDeployInfo{
					Id:         deploy.Id,
					StartTime:  startTime.Format("15:04:05"),
					EndTime:    endTime.Format("15:04:05"),
					Note:       deploy.Note,
					CourseSort: deploy.CourseSort,
					Grouping:   deploy.Grouping,
					CourseFlag: deploy.CourseFlag,
				})
			}
		}
	}

	return &res, nil
}
