package logic

import (
	"context"
	"demacia/cloudscreen/courserecord/cron/internal/config"
	"demacia/cloudscreen/courserecord/model"
	"demacia/datacenter/class/rpc/classclient"
	"demacia/datacenter/coursetable/rpc/coursetableclient"
	"demacia/datacenter/member/rpc/memberclient"
	"fmt"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
	"time"
)

type (
	RecordGenerateLogic struct {
		logx.Logger
		CourseRecordModel       *model.CourseRecordModel
		CourseRecordConfigModel *model.CourseRecordConfigModel
		CourseRecordCountModel  *model.CourseRecordCountModel
		CourseRecordDateModel   *model.CourseRecordDateModel
		CourseTableRpc          coursetableclient.Coursetable
		ClassRpc                classclient.Class
		MemberRpc               memberclient.Member
	}
)

func NewRecordGenerateLogic(ctx context.Context, c *config.Config) *RecordGenerateLogic {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &RecordGenerateLogic{
		Logger:                  logx.WithContext(ctx),
		CourseRecordModel:       model.NewCourseRecordModel(conn, cacheRedis),
		CourseRecordConfigModel: model.NewCourseRecordConfigModel(conn, cacheRedis),
		CourseRecordCountModel:  model.NewCourseRecordCountModel(conn, cacheRedis),
		CourseRecordDateModel:   model.NewCourseRecordDateModel(conn, cacheRedis),
		CourseTableRpc:          coursetableclient.NewCoursetable(zrpc.MustNewClient(c.CourseTableRpc)),
		MemberRpc:               memberclient.NewMember(zrpc.MustNewClient(c.MemberRpc)),
		ClassRpc:                classclient.NewClass(zrpc.MustNewClient(c.ClassRpc)),
	}
}

func (l *RecordGenerateLogic) GenerateCourseRecord() {

	configs, err := l.CourseRecordConfigModel.SelectList()
	if err != nil {
		l.Logger.Errorf("load course_record_config err:%s", err.Error())
		return
	}
	var orgIds []int64
	classIds := map[int64][]int64{}
	classCourseTableInfo := map[int64]*coursetableclient.CourseTableRecordResp{}
	configMap := map[int64]*model.CourseRecordConfig{}
	students := map[int64][]struct {
		id   int64
		name string
	}{}
	for _, c := range configs {
		if c.SignHoliday {
			val, err := l.CourseRecordDateModel.SelectTodayIsHoliday()
			if err != nil {
				l.Logger.Errorf("select course_record_date err:%s", err.Error())
			}
			if val > 0 {
				continue
			}
		}
		configMap[c.OrgId] = c
		orgIds = append(orgIds, c.OrgId)
	}
	if len(orgIds) == 0 {
		return
	}
	for _, oid := range orgIds {
		classes, err := l.ClassRpc.ListByOrgId(context.Background(), &classclient.OrgIdReq{OrgId: oid})
		if err != nil {
			continue
		}
		for _, class := range classes.List {
			classIds[oid] = append(classIds[oid], class.Id)
		}
	}
	if len(classIds) == 0 {
		return
	}

	for k, v := range classIds {
		for _, cid := range v {
			deploys, err := l.CourseTableRpc.GetCourseTableRecordInfo(context.Background(), &coursetableclient.OrgIdAndClassIdReq{OrgId: k, ClassId: cid})
			if err != nil {
				continue
			}
			if len(deploys.List) == 0 {
				continue
			}
			classCourseTableInfo[cid] = deploys
		}
	}
	for k := range classCourseTableInfo {
		for i := 1; i <= 3; i++ {
			students[k] = append(students[k], struct {
				id   int64
				name string
			}{id: int64(i), name: fmt.Sprintf("%s%d", "学生", i)})
		}
	}

	for k, v := range students {
		for _, student := range v {
			var records model.CourseRecords
			for _, course := range classCourseTableInfo[k].List {
				if configMap[course.OrgId].SignPerson == 1 || configMap[course.OrgId].SignPerson == 3 {
					records = append(records, &model.CourseRecord{
						OrgId:        course.OrgId,
						UserId:       student.id,
						Truename:     student.name,
						UserType:     1,
						SignDate:     time.Now().Format("2006-01-02"),
						SignTime:     "00:00:00",
						SubjectName:  course.SubjectName,
						CourseNote:   course.CourseNote,
						ClassName:    course.ClassName,
						ClassId:      k,
						PositionName: course.PositionName,
						StartTime:    course.StartTime,
						EndTime:      course.EndTime,
					})
				} else {
					break
				}
			}
			if err := l.CourseRecordModel.InsertCourseRecords(records); err != nil {
				l.Logger.Errorf("insert course_record error:%s,student_id:%d", err.Error(), student.id)
				continue
			}
			if err := l.CourseRecordCountModel.InsertCourseRecordCount(&model.CourseRecordCount{
				UserId:      student.id,
				Truename:    student.name,
				UserType:    1,
				ShouldCount: int64(len(records)),
				CountDate:   time.Now().Format("2006-01-02"),
				ClassId:     k,
			}); err != nil {
				l.Logger.Errorf("insert courses_count error:%s,student_id:%d", err.Error(), student.id)
				continue
			}
		}
	}

	type mInfo struct {
		name  string
		count int64
	}

	memberCount := map[int64]*mInfo{}

	for k, v := range classCourseTableInfo {
		var records model.CourseRecords
		for _, course := range v.List {
			if configMap[course.OrgId].SignPerson > 2 {
				if _, ok := memberCount[course.MemberId]; !ok {
					memberCount[course.MemberId] = &mInfo{
						name:  course.MemberName,
						count: 1,
					}
				} else {
					memberCount[course.MemberId].count++
				}
				records = append(records, &model.CourseRecord{
					OrgId:        course.OrgId,
					UserId:       course.MemberId,
					Truename:     course.MemberName,
					UserType:     2,
					SignDate:     time.Now().Format("2006-01-02"),
					SignTime:     "00:00:00",
					SubjectName:  course.SubjectName,
					CourseNote:   course.CourseNote,
					ClassName:    course.ClassName,
					ClassId:      k,
					PositionName: course.PositionName,
					StartTime:    course.StartTime,
					EndTime:      course.EndTime,
				})
			} else {
				break
			}
		}
		if len(records) > 0 {
			if err := l.CourseRecordModel.InsertCourseRecords(records); err != nil {
				l.Logger.Errorf("insert course_record error:%s,member_id:%d", err.Error(), records[0].UserId)
				continue
			}
		}
	}

	for id, info := range memberCount {
		if err := l.CourseRecordCountModel.InsertCourseRecordCount(&model.CourseRecordCount{
			UserId:      id,
			Truename:    info.name,
			UserType:    2,
			ShouldCount: info.count,
			CountDate:   time.Now().Format("2006-01-02"),
		}); err != nil {
			l.Logger.Errorf("insert courses_count error:%s,member_id:%d", err.Error(), id)
			continue
		}
	}

}
