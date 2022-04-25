package logic

import (
	"context"
	"demacia/cloudscreen/courserecord/api/internal/svc"
	"demacia/cloudscreen/courserecord/api/internal/types"
	"demacia/cloudscreen/courserecord/errors"
	"demacia/cloudscreen/courserecord/model"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/common/file"
	"encoding/base64"
	"github.com/google/uuid"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type CourseRecordAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseRecordAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) CourseRecordAddLogic {
	return CourseRecordAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseRecordAddLogic) CourseRecordAdd(req types.CourseRecordAddReq) (*types.SuccessReply, error) {
	resp := &types.SuccessReply{Success: false}
	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, errlist.NoAuth
	}
	signTime, err := time.Parse("15:04:05", req.SignTime)
	if err != nil {
		return resp, errlist.InvalidParam
	}
	records, err := l.svcCtx.CourseRecordModel.SelectByUserIdAndDate(req.UserId, time.Now().Format("2006-01-02"))
	if err != nil || len(records) == 0 {
		return resp, errors.SignRecordNotFound
	}
	for _, r := range records {
		startTime, _ := time.Parse("2006-01-02T15:04:05Z", r.StartTime)
		endTime, _ := time.Parse("2006-01-02T15:04:05Z", r.EndTime)
		if signTime.Add(10*time.Minute).After(startTime) && signTime.Add(10*time.Minute).Before(endTime) {
			if r.Status != 0 {
				return resp, errors.SignRecordExist
			}
			if signTime.Before(startTime) {
				r.Status = 1
			} else {
				r.Status = 2
			}
			r.SignTime = signTime.Format("15:04:05")
			if err := l.svcCtx.CourseRecordModel.UpdateCourseRecord(r); err != nil {
				return resp, errlist.Unknown
			}
			recordCount, err := l.svcCtx.CourseRecordCountModel.SelectByUserIdToToday(r.UserId)
			if err != nil {
				return resp, errlist.Unknown
			}
			if r.Status == 1 {
				recordCount.NormalCount++
			} else if r.Status == 2 {
				recordCount.LateCount++
			}
			if err := l.svcCtx.CourseRecordCountModel.UpdateRecordCount(&recordCount); err != nil {
				return resp, errlist.Unknown
			}
			threading.GoSafe(func() {
				req.Photo = strings.Replace(req.Photo, `\n`, "", -1)
				filename := uuid.New().String() + ".png"
				ddd, err := base64.StdEncoding.DecodeString(req.Photo)
				if err != nil {
					logx.Errorf("base64 err:%s", err.Error())
					return
				}
				imagePath := os.TempDir() + filename
				err = ioutil.WriteFile(imagePath, ddd, 0666)
				if err != nil {
					logx.Errorf("write base64 into file err:%s", err.Error())
					return
				}
				remotePath, err := file.Upload(&file.UpdateReq{
					Cache:     l.svcCtx.Config.CacheRedis.NewRedis(),
					LocalPath: imagePath,
					FileName:  filename,
					OrgId:     oid,
					IsTmp:     false,
				})
				if err != nil {
					logx.Errorf("upload image to oss err:%s", err.Error())
					return
				}
				_ = l.svcCtx.CourseRecordModel.UpdatePhoto(&model.CourseRecord{
					Id:    r.Id,
					Photo: remotePath,
				})
				os.Remove(imagePath)
			})
			break
		}
	}
	resp.Success = true
	return resp, nil
}
