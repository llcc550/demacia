package logic

import (
	"context"
	"demacia/datacenter/coursetable/errors"
	"fmt"
	"time"

	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GenerateCourseTableDeployLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateCourseTableDeployLogic(ctx context.Context, svcCtx *svc.ServiceContext) GenerateCourseTableDeployLogic {
	return GenerateCourseTableDeployLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateCourseTableDeployLogic) GenerateCourseTableDeploy(req types.DeployGenerateReq) (*types.DeployGenerateReply, error) {

	var deploys types.DeployGenerateReply
	resp := &types.DeployGenerateReply{CourseTableDeploy: []*types.CourseTableGenerate{}}
	morningStartTime, err := time.Parse("15:04:05", req.MorningStartTime)
	if err != nil {
		return resp, errors.InvalidTime
	}
	afternoonStartTime, err := time.Parse("15:04:05", req.AfternoonStartTime)
	if err != nil {
		return resp, errors.InvalidTime
	}

	morningSelfStudyStartTime, err := time.Parse("15:04:05", req.MorningSelfStudyStartTime)
	if err != nil {
		return resp, errors.InvalidTime
	}

	nightSelfStudyStartTime, err := time.Parse("15:04:05", req.NightSelfStudyStartTime)
	if err != nil {
		return resp, errors.InvalidTime
	}

	courseSort := int8(1)
	if req.MorningSelfStudyCount > 0 {
		var minuteCount int64
		for i := int8(1); i <= req.MorningSelfStudyCount; i++ {
			if i == 1 {
				deploys.CourseTableDeploy = append(deploys.CourseTableDeploy, &types.CourseTableGenerate{
					StartTime:  morningSelfStudyStartTime.Format("15:04:05"),
					EndTime:    morningSelfStudyStartTime.Add(time.Duration(req.MorningSelfStudyTime) * time.Minute).Format("15:04:05"),
					Note:       fmt.Sprintf("早自习第%d节", i),
					Grouping:   1,
					CourseSort: courseSort,
					CourseFlag: 1,
				})
				courseSort++
			} else {
				deploys.CourseTableDeploy = append(deploys.CourseTableDeploy, &types.CourseTableGenerate{
					StartTime:  morningSelfStudyStartTime.Add(time.Duration(minuteCount) * time.Minute).Format("15:04:05"),
					EndTime:    morningSelfStudyStartTime.Add(time.Duration(minuteCount+int64(req.MorningCourseTime)) * time.Minute).Format("15:04:05"),
					Note:       fmt.Sprintf("早自习第%d节", i),
					Grouping:   1,
					CourseSort: courseSort,
					CourseFlag: 1,
				})
				courseSort++
			}
			minuteCount += int64(req.MorningSelfStudyTime + req.BreakTime)
		}
		if morningSelfStudyStartTime.Add(time.Duration(minuteCount) * time.Minute).After(morningStartTime) {
			return resp, errors.InvalidMorningSelfTime
		}
	}

	var minuteCount int64
	courseIndex := 1
	for i := 1; i <= int(req.MorningCourseCount); i++ {
		info := types.CourseTableGenerate{
			CourseSort: courseSort,
			Grouping:   2,
		}
		if i == 1 {
			info.StartTime = morningStartTime.Format("15:04:05")
			info.EndTime = morningStartTime.Add(time.Duration(req.MorningCourseTime) * time.Minute).Format("15:04:05")
			info.Note = fmt.Sprintf("%d节", courseIndex)
			courseIndex++
			minuteCount += int64(req.MorningCourseTime + req.BreakTime)
			info.CourseFlag = 1
		} else if int8(i) == req.MorningRecessIndex {
			info.StartTime = morningStartTime.Add(time.Duration(minuteCount-int64(req.BreakTime)) * time.Minute).Format("15:04:05")
			minuteCount -= int64(req.BreakTime)
			info.EndTime = morningStartTime.Add(time.Duration(minuteCount+int64(req.MorningRecessTime)) * time.Minute).Format("15:04:05")
			info.Note = "上午大课间"
			info.Grouping = 3
			minuteCount += int64(req.MorningRecessTime)
		} else {
			info.StartTime = morningStartTime.Add(time.Duration(minuteCount) * time.Minute).Format("15:04:05")
			info.EndTime = morningStartTime.Add(time.Duration(minuteCount+int64(req.MorningCourseTime)) * time.Minute).Format("15:04:05")
			info.Note = fmt.Sprintf("%d节", courseIndex)
			courseIndex++
			minuteCount += int64(req.MorningCourseTime + req.BreakTime)
			info.CourseFlag = 1
		}
		deploys.CourseTableDeploy = append(deploys.CourseTableDeploy, &info)
		courseSort++
	}

	mVal, err := time.Parse("15:04:05", deploys.CourseTableDeploy[len(deploys.CourseTableDeploy)-1].EndTime)
	if mVal.After(afternoonStartTime) {
		return resp, errors.CourseTimeErr
	}

	deploys.CourseTableDeploy = append(deploys.CourseTableDeploy, &types.CourseTableGenerate{
		StartTime:  deploys.CourseTableDeploy[len(deploys.CourseTableDeploy)-1].EndTime,
		EndTime:    req.AfternoonStartTime,
		Note:       "午休",
		CourseSort: courseSort,
		CourseFlag: 0,
		Grouping:   4,
	})
	courseSort++
	minuteCount = 0
	for i := req.MorningCourseCount; i < req.AfternoonCourseCount+req.MorningCourseCount; i++ {
		info := types.CourseTableGenerate{
			CourseSort: courseSort,
			Grouping:   5,
		}
		if i == req.MorningCourseCount {
			info.StartTime = afternoonStartTime.Format("15:04:05")
			info.EndTime = afternoonStartTime.Add(time.Duration(req.AfternoonCourseTime) * time.Minute).Format("15:04:05")
			info.Note = fmt.Sprintf("%d节", courseIndex)
			courseIndex++
			minuteCount += int64(req.AfternoonCourseTime + req.BreakTime)
			info.CourseFlag = 1
		} else if i+1 == req.AfternoonRecessIndex {
			info.StartTime = afternoonStartTime.Add(time.Duration(minuteCount-int64(req.BreakTime)) * time.Minute).Format("15:04:05")
			minuteCount -= int64(req.BreakTime)
			info.EndTime = afternoonStartTime.Add(time.Duration(minuteCount+int64(req.AfternoonRecessTime)) * time.Minute).Format("15:04:05")
			info.Note = "下午大课间"
			minuteCount += int64(req.AfternoonRecessTime)
			info.Grouping = 6
		} else {
			info.CourseFlag = 1
			info.StartTime = afternoonStartTime.Add(time.Duration(minuteCount) * time.Minute).Format("15:04:05")
			info.EndTime = afternoonStartTime.Add(time.Duration(minuteCount+int64(req.AfternoonCourseTime)) * time.Minute).Format("15:04:05")
			info.Note = fmt.Sprintf("%d节", courseIndex)
			courseIndex++
			minuteCount += int64(req.AfternoonCourseTime + req.BreakTime)
		}
		deploys.CourseTableDeploy = append(deploys.CourseTableDeploy, &info)
		courseSort++
	}

	nightSelfTimeVal, _ := time.Parse("15:04:05", deploys.CourseTableDeploy[len(deploys.CourseTableDeploy)-1].EndTime)

	if req.NightSelfStudyCount > 0 && nightSelfTimeVal.After(nightSelfStudyStartTime) {
		return resp, errors.InvalidNightSelfTime
	}

	if req.NightSelfStudyCount > 0 {
		var minuteCount int64
		for i := int8(1); i <= req.NightSelfStudyCount; i++ {
			if i == 1 {
				deploys.CourseTableDeploy = append(deploys.CourseTableDeploy, &types.CourseTableGenerate{
					StartTime:  nightSelfStudyStartTime.Format("15:04:05"),
					EndTime:    nightSelfStudyStartTime.Add(time.Duration(req.MorningSelfStudyTime) * time.Minute).Format("15:04:05"),
					Note:       fmt.Sprintf("晚自习第%d节", i),
					Grouping:   7,
					CourseSort: courseSort,
					CourseFlag: 1,
				})
				courseSort++
			} else {
				deploys.CourseTableDeploy = append(deploys.CourseTableDeploy, &types.CourseTableGenerate{
					StartTime:  nightSelfStudyStartTime.Add(time.Duration(minuteCount) * time.Minute).Format("15:04:05"),
					EndTime:    nightSelfStudyStartTime.Add(time.Duration(minuteCount+int64(req.MorningCourseTime)) * time.Minute).Format("15:04:05"),
					Grouping:   7,
					Note:       fmt.Sprintf("晚自习第%d节", i),
					CourseSort: courseSort,
					CourseFlag: 1,
				})
				courseSort++
			}
			minuteCount += int64(req.MorningSelfStudyTime + req.BreakTime)
		}
	}

	for i := int8(1); i <= 7; i++ {
		for _, deploy := range deploys.CourseTableDeploy {
			resp.CourseTableDeploy = append(resp.CourseTableDeploy, &types.CourseTableGenerate{
				StartTime:  deploy.StartTime,
				EndTime:    deploy.EndTime,
				Note:       deploy.Note,
				Weekday:    i,
				CourseSort: deploy.CourseSort,
				Grouping:   deploy.Grouping,
				CourseFlag: deploy.CourseFlag,
			})
		}
	}

	return resp, nil
}
