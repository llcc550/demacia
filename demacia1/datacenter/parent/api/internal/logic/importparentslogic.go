package logic

import (
	"context"
	"demacia/common/datacenter"
	"fmt"
	"os"
	"strings"

	"demacia/common/baseauth"
	"demacia/common/basefunc"
	"demacia/common/errlist"
	"demacia/datacenter/parent/api/internal/svc"
	"demacia/datacenter/parent/api/internal/types"
	"demacia/datacenter/parent/model"
	"demacia/datacenter/student/rpc/student"
	"demacia/service/websocket/rpc/websocket"

	"github.com/tealeg/xlsx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

type ImportParentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImportParentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ImportParentsLogic {
	return ImportParentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var relationMap = map[string]int8{
	"家长": 0,
	"爸爸": 1,
	"妈妈": 2,
	"爷爷": 3,
	"奶奶": 4,
	"外公": 5,
	"外婆": 6,
}

type StudentInfo struct {
	StudentId   int64  `json:"student_id"`
	StudentName string `json:"student_name"`
	ClassId     int64  `json:"class_id"`
	ClassName   string `json:"class_name"`
	Relation    int8   `json:"relation"`
}

func (l *ImportParentsLogic) ImportParents(req types.UrlReq) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	threading.GoSafe(func() {
		l.importParents(orgId, req.Url, req.WebsocketUuid)
	})
	return nil
}

func (l *ImportParentsLogic) importParents(orgId int64, url, uuid string) {
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
	for _, sheet := range xlFile.Sheets {
		for k, row := range sheet.Rows {
			if k == 0 {
				continue
			}
			var parentId int64
			parentName := strings.TrimSpace(row.Cells[0].Value)
			mobile := strings.TrimSpace(row.Cells[1].Value)
			address := strings.TrimSpace(row.Cells[2].Value)
			idNumber := strings.TrimSpace(row.Cells[3].Value)
			if idNumber != "" {
				state := basefunc.CheckIdCard(idNumber)
				if !state {
					fn(errlist.ImportErr, idNumber, "身份证号格式有误")
				}
				continue
			}
			parentInfo, err := l.svcCtx.ParentModel.FindOneByMobile(mobile)
			if err == nil {
				parentId = parentInfo.Id
			} else {
				parentId, err = l.svcCtx.ParentModel.InsertOne(&model.Parent{
					TrueName:   parentName,
					UserName:   mobile,
					Password:   basefunc.HashPassword(mobile),
					Face:       "",
					FaceStatus: 0,
					IdNumber:   idNumber,
					Address:    address,
					Mobile:     mobile,
					Pinyin:     basefunc.GetPinyin(parentName),
				})
				if err != nil {
					fn(errlist.ImportErr, "导入家长失败")
					continue
				}
			}
			studentUserName := strings.TrimSpace(row.Cells[4].Value)
			studentTrueName := strings.TrimSpace(row.Cells[5].Value)
			relation := strings.TrimSpace(row.Cells[6].Value)
			if parentName == "" || mobile == "" || studentUserName == "" || relation == "" {
				break
			}

			if _, ok := relationMap[relation]; !ok {
				fn(errlist.ImportErr, "关系错误")
				continue
			}
			studentMap := map[int64]StudentInfo{}
			studentInfo, err := l.svcCtx.StudentRpc.FindOneByUserName(context.Background(), &student.UserNameRequest{
				OrgId:    orgId,
				UserName: studentUserName,
			})
			if err != nil || studentInfo.Id == 0 {
				fmt.Println(err)
				fn(errlist.ImportErr, studentUserName, "学生不存在")
				continue
			}
			studentMap[studentInfo.Id] = StudentInfo{
				StudentId:   studentInfo.Id,
				StudentName: studentInfo.Name,
				ClassId:     studentInfo.ClassId,
				ClassName:   studentInfo.ClassName,
				Relation:    relationMap[relation],
			}
			if len(row.Cells) > 7 && len(row.Cells) < 11 {
				studentUserNameTwo := strings.TrimSpace(row.Cells[7].Value)
				studentTrueNameTwo := strings.TrimSpace(row.Cells[8].Value)
				relationTwo := strings.TrimSpace(row.Cells[9].Value)
				if studentUserNameTwo != "" {
					studentInfo1, err := l.svcCtx.StudentRpc.FindOneByUserName(context.Background(), &student.UserNameRequest{
						OrgId:    orgId,
						UserName: studentUserNameTwo,
					})
					if err != nil || studentInfo1.Id == 0 {
						fn(errlist.ImportErr, studentUserNameTwo, "学生不存在")
						continue
					}
					studentMap[studentInfo1.Id] = StudentInfo{
						StudentId:   studentInfo1.Id,
						StudentName: studentTrueNameTwo,
						ClassId:     studentInfo1.ClassId,
						ClassName:   studentInfo1.ClassName,
						Relation:    relationMap[relationTwo],
					}
				}
			}
			if len(row.Cells) > 10 && len(row.Cells) < 14 {
				studentUserNameTwo := strings.TrimSpace(row.Cells[7].Value)
				studentTrueNameTwo := strings.TrimSpace(row.Cells[8].Value)
				relationTwo := strings.TrimSpace(row.Cells[9].Value)
				if studentUserNameTwo != "" {
					studentInfo1, err := l.svcCtx.StudentRpc.FindOneByUserName(context.Background(), &student.UserNameRequest{
						OrgId:    orgId,
						UserName: studentUserNameTwo,
					})
					if err != nil || studentInfo1.Id == 0 {
						fn(errlist.ImportErr, studentUserNameTwo, "学生不存在")
						continue
					}
					studentMap[studentInfo1.Id] = StudentInfo{
						StudentId:   studentInfo1.Id,
						StudentName: studentTrueNameTwo,
						ClassId:     studentInfo1.ClassId,
						ClassName:   studentInfo1.ClassName,
						Relation:    relationMap[relationTwo],
					}
				}
				studentUserNameThree := strings.TrimSpace(row.Cells[10].Value)
				studentTrueNameThree := strings.TrimSpace(row.Cells[11].Value)
				relationThree := strings.TrimSpace(row.Cells[12].Value)
				if studentUserNameThree != "" {
					studentInfo2, err := l.svcCtx.StudentRpc.FindOneByUserName(context.Background(), &student.UserNameRequest{
						OrgId:    orgId,
						UserName: studentUserNameThree,
					})
					if err != nil || studentInfo2.Id == 0 {
						fn(errlist.ImportErr, studentUserNameThree, "学生不存在")
						continue
					}
					studentMap[studentInfo2.Id] = StudentInfo{
						StudentId:   studentInfo2.Id,
						StudentName: studentTrueNameThree,
						ClassId:     studentInfo2.ClassId,
						ClassName:   studentInfo2.ClassName,
						Relation:    relationMap[relationThree],
					}
				}
			}
			datacenterStudentInfoSlice := make([]datacenter.StudentInfo, 0)
			for _, studentInfo := range studentMap {
				if studentInfo.StudentId == 0 {
					continue
				}
				datacenterStudentInfoSlice = append(datacenterStudentInfoSlice, datacenter.StudentInfo{
					StudentId: studentInfo.StudentId,
					Relation:  studentInfo.Relation,
				})
				relationInfo, err := l.svcCtx.StudentParentModel.FindOneByParentIdAndStudentId(parentId, studentInfo.StudentId)
				if err != nil {
					_, err := l.svcCtx.StudentParentModel.InsertOne(&model.StudentParent{
						OrgId:       orgId,
						ClassId:     studentInfo.ClassId,
						ClassName:   studentInfo.ClassName,
						ParentId:    parentId,
						StudentId:   studentInfo.StudentId,
						StudentName: studentInfo.StudentName,
						Relation:    studentInfo.Relation,
					})
					if err != nil {
						fn(errlist.ImportSuc, studentUserName+studentTrueName)
						continue
					}
				} else {
					err := l.svcCtx.StudentParentModel.UpdateOne(&model.StudentParent{
						Id:          relationInfo.Id,
						OrgId:       orgId,
						ClassId:     studentInfo.ClassId,
						ClassName:   studentInfo.ClassName,
						ParentId:    parentId,
						StudentId:   studentInfo.StudentId,
						StudentName: studentInfo.StudentName,
						Relation:    studentInfo.Relation,
					})
					if err != nil {
						continue
					}
				}

			}
			s := datacenter.Marshal(datacenter.Parent, parentId, datacenter.Add, datacenter.ParentData{
				ParentName:  parentName,
				UserName:    mobile,
				Mobile:      mobile,
				StudentInfo: datacenterStudentInfoSlice,
			})
			_ = l.svcCtx.KqPusher.Push(s)
		}
	}
}
