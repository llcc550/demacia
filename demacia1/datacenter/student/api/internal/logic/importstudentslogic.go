package logic

import (
	"context"
	"fmt"
	"os"
	"strings"

	"demacia/common/baseauth"
	"demacia/common/basefunc"
	"demacia/common/datacenter"
	"demacia/common/errlist"
	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/student/api/internal/svc"
	"demacia/datacenter/student/api/internal/types"
	"demacia/datacenter/student/model"
	"demacia/service/websocket/rpc/websocket"

	"github.com/tealeg/xlsx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

type ImportStudentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImportStudentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ImportStudentsLogic {
	return ImportStudentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImportStudentsLogic) ImportStudents(req types.UrlReq) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	threading.GoSafe(func() {
		l.importStudents(orgId, req.Url, req.WebsocketUuid)
	})
	return nil
}

func (l *ImportStudentsLogic) importStudents(orgId int64, url, uuid string) {
	fn := func(err *httpx.CodeError, msg ...interface{}) {
		_, _ = l.svcCtx.WebsocketRpc.Push(context.Background(), &websocket.Request{
			Key:  uuid,
			Code: int64(err.Code()),
			Msg:  fmt.Sprintf(err.Error(), msg...),
		})
	}
	path := basefunc.Download(url)
	if path == "" {
		fn(errlist.Excel)
		return
	}
	defer os.Remove(path)
	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		fn(errlist.Excel)
		return
	}
	classIds := make([]int64, 0)
	classMap := map[string]*class.ClassInfo{}
	for _, sheet := range xlFile.Sheets {
		for k, row := range sheet.Rows {
			if k == 0 {
				continue
			}
			if len(row.Cells) < 3 {
				break
			}
			classTitle := strings.TrimSpace(row.Cells[0].Value)
			studentUserName := strings.TrimSpace(row.Cells[1].Value)
			studentTrueName := strings.TrimSpace(row.Cells[2].Value)
			studentSex := strings.TrimSpace(row.Cells[3].Value)
			if classTitle == "" || studentUserName == "" || studentTrueName == "" {
				break
			}
			sex := int8(0)
			if studentSex == "男" {
				sex = 1
			}
			if _, ok := classMap[classTitle]; !ok {
				classInfo, err := l.svcCtx.ClassRpc.GetClassInfoByFullName(context.Background(), &class.FullNameReq{
					OrgId:    orgId,
					FullName: classTitle,
				})
				if err != nil || classInfo.Id == 0 {
					fn(errlist.ImportErr, classTitle+studentTrueName, "班级不存在")
					continue
				}
				classMap[classTitle] = classInfo
			}
			classInfo := classMap[classTitle]
			classIds = append(classIds, classInfo.Id)
			studentInfo, err := l.svcCtx.StudentModel.FindOneByOrgIdAndUserName(orgId, studentUserName)
			if err == nil {
				if studentInfo.ClassId != classInfo.Id {
					classIds = append(classIds, studentInfo.ClassId)
				}
				studentId := studentInfo.Id
				err := l.svcCtx.StudentModel.UpdateOne(&model.Student{
					Id:            studentId,
					OrgId:         orgId,
					ClassId:       classInfo.Id,
					TrueName:      studentTrueName,
					UserName:      studentUserName,
					Deleted:       false,
					StageId:       classInfo.StageId,
					GradeId:       classInfo.GradeId,
					ClassFullName: classTitle,
					Password:      basefunc.HashPassword(studentUserName),
					Sex:           sex,
				})
				if err != nil {
					fn(errlist.ImportSuc, classTitle+studentTrueName)
					continue
				}
				s := datacenter.Marshal(datacenter.Student, studentId, datacenter.Update, datacenter.StudentData{
					StudentName:    studentTrueName,
					UserName:       studentUserName,
					OrganizationId: orgId,
					ClassId:        classInfo.Id,
					Sex:            sex,
				})
				_ = l.svcCtx.KqPusher.Push(s)
			} else {
				studentId, err := l.svcCtx.StudentModel.InsertOne(&model.Student{
					OrgId:         orgId,
					ClassId:       classInfo.Id,
					TrueName:      studentTrueName,
					UserName:      studentUserName,
					Deleted:       false,
					StageId:       classInfo.StageId,
					GradeId:       classInfo.GradeId,
					ClassFullName: classTitle,
					Password:      basefunc.HashPassword(studentUserName),
					Sex:           sex,
				})
				if err != nil {
					fn(errlist.ImportSuc, classTitle+studentTrueName)
					continue
				}
				s := datacenter.Marshal(datacenter.Student, studentId, datacenter.Add, datacenter.StudentData{
					StudentName:    studentTrueName,
					UserName:       studentUserName,
					OrganizationId: orgId,
					ClassId:        classInfo.Id,
					Sex:            sex,
				})
				_ = l.svcCtx.KqPusher.Push(s)
			}
		}
	}

	for _, classId := range classIds {
		count, err := l.svcCtx.StudentModel.GetClassStudentCount(classId)
		if err != nil {
			continue
		}
		_, _ = l.svcCtx.ClassRpc.ChangeStudentNum(context.Background(), &class.StudentNumReq{
			ClassId:    classId,
			StudentNum: count,
		})
	}
}
