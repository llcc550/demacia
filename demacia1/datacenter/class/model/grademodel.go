package model

import (
	"demacia/common/basefunc"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"demacia/common/cachemodel"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"xorm.io/builder"
)

var (
	gradeFieldNames         = strings.Join(cachemodel.RawFieldNames(&Grade{}, true), ",")
	cacheClassGradeIdPrefix = "cache:class:grade:id:"
)

type (
	GradeModel struct {
		*cachemodel.CachedModel
	}

	Grade struct {
		Id             int64  `db:"id"`
		OrgId          int64  `db:"org_id"`           // 所属机构
		Title          string `db:"title"`            // 年级名称
		StageId        int64  `db:"stage_id"`         // 学段ID。冗余数据
		StageTitle     string `db:"stage_title"`      // 学段名称。冗余数据
		GradeMemberNum int64  `db:"grade_member_num"` // 年级人数
	}
	Grades []*Grade
)

func NewGradeModel(conn sqlx.SqlConn, cache *redis.Redis) *GradeModel {
	return &GradeModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"class"."grade"`, cache),
	}
}

func (m *GradeModel) ListByOrgIdAndStageId(orgId, stageId int64) (Grades, error) {
	sql, args, _ := builder.Postgres().
		Select(gradeFieldNames).
		From(m.Table).
		Where(builder.Eq{"org_id": orgId, "stage_id": stageId, "deleted": "false"}).
		ToSQL()
	var res Grades
	err := m.Conn.QueryRows(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *GradeModel) Insert(orgId, stageId int64, gradeTitle, stageTitle string) (int64, error) {
	sqlString, args, _ := builder.Postgres().Insert(builder.Eq{
		"org_id":      orgId,
		"title":       gradeTitle,
		"stage_id":    stageId,
		"stage_title": stageTitle,
	}).Into(m.Table).ToSQL()
	var LastInsertId int64
	err := m.Conn.QueryRow(&LastInsertId, sqlString+" returning id", args...)
	return LastInsertId, err
}

func (m *GradeModel) Update(orgId, Id int64, gradeTitle string) error {
	sqlString, args, _ := builder.Postgres().Update(builder.Eq{
		"title": gradeTitle,
	}).From(m.Table).Where(builder.Eq{"id": Id, "org_id": orgId, "deleted": "false"}).ToSQL()
	_, err := m.ExecNoCache(sqlString, args...)
	if err != nil {
		return err
	}
	return nil

}

func (m *GradeModel) GetGradeByIdAndStageId(Id, stageId int64) (*Grade, error) {
	sql, args, _ := builder.Postgres().
		Select(gradeFieldNames).
		From(m.Table).
		Where(builder.Eq{"id": Id, "stage_id": stageId, "deleted": "false"}).
		ToSQL()
	var res Grade
	err := m.Conn.QueryRow(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *GradeModel) FindGradeById(Id int64) (*Grade, error) {
	sql, args, _ := builder.Postgres().
		Select(gradeFieldNames).
		From(m.Table).
		Where(builder.Eq{"id": Id, "deleted": "false"}).
		ToSQL()
	var res Grade
	err := m.Conn.QueryRow(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (m *GradeModel) GetStageYearByStageId(stageId int64) (int, error) {
	var res int
	sql, args, _ := builder.Postgres().
		Select("COUNT(id)").
		From(m.Table).
		Where(builder.Eq{"stage_id": stageId, "deleted": false}).
		ToSQL()
	err := m.Conn.QueryRow(&res, sql, args...)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (m *GradeModel) GetGradeListByOrgId(Id int64) (*Grades, error) {
	sql, args, _ := builder.Postgres().
		Select(gradeFieldNames).
		From(m.Table).
		Where(builder.Eq{"org_id": Id, "deleted": "false"}).
		ToSQL()
	var res Grades
	err := m.Conn.QueryRows(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (m *GradeModel) GetGradeListWithId(Id []int64) (*Grades, error) {
	sql, args, _ := builder.Postgres().
		Select(gradeFieldNames).
		From(m.Table).
		Where(builder.In("id", Id)).
		And(builder.Eq{"deleted": "false"}).
		ToSQL()
	var res Grades
	err := m.Conn.QueryRows(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *GradeModel) DeleteGradeById(orgId, Id int64) error {
	sqlString, args, _ := builder.Postgres().Update(builder.Eq{
		"deleted": "true",
	}).From(m.Table).Where(builder.Eq{"id": Id, "org_id": orgId}).ToSQL()
	_, err := m.ExecNoCache(sqlString, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m *GradeModel) ChangeStudentNumById(GradeId, Num int64) error {
	sqlString, args, _ := builder.Postgres().Update(builder.Eq{
		"grade_member_num": Num,
	}).Where(builder.Eq{
		"id":      GradeId,
		"deleted": false,
	}).
		From(m.Table).ToSQL()
	_, err := m.ExecNoCache(sqlString, args...)
	return err
}

// ============================== stage 根据学段更新数据=====================================

func (m *GradeModel) UpdateStageTitleByStageId(orgId, stageId int64, stageTitle string) error {
	sqlString, args, _ := builder.Postgres().Update(builder.Eq{
		"stage_title": stageTitle,
	}).From(m.Table).Where(builder.Eq{"stage_id": stageId, "org_id": orgId, "deleted": "false"}).ToSQL()
	_, err := m.ExecNoCache(sqlString, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m *GradeModel) DeleteGradeByStageId(orgId, stageId int64) error {
	sqlString, args, _ := builder.Postgres().Update(builder.Eq{
		"deleted": "true",
	}).From(m.Table).Where(builder.Eq{"stage_id": stageId, "org_id": orgId}).ToSQL()
	_, err := m.ExecNoCache(sqlString, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m *GradeModel) AddGradeNumByStageId(orgId, stageId int64, stageTitle string, num int64) error {
	grades := make(Grades, 0)
	for i := 1; i <= int(num); i++ {
		ca, err := IntToCa(i)
		if err != nil {
			ca = strconv.Itoa(i)
		}
		grades = append(grades, &Grade{
			Title:          ca + "年级",
			OrgId:          orgId,
			StageId:        stageId,
			StageTitle:     stageTitle,
			GradeMemberNum: 0,
		})
	}
	query, args, err := basefunc.BatchInsertString(m.Table, grades)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = m.ExecNoCache(query, args...)
	return err
}

func IntToCa(n int) (string, error) {
	nStr := strconv.Itoa(n)
	var company = []string{"", "十", "百", "千", "万"}
	var zhCa = []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}

	var res string
	if len(nStr) > len(company) {
		return "", errors.New("the length of parameter 1 is out of range")
	}
	zero := false
	for i := 1; i <= len(nStr); i++ {
		site, _ := strconv.Atoi(nStr[i-1 : i])
		if i == len(nStr) && site == 0 {
			break
		}
		if site == 0 {
			if !zero {
				res += zhCa[site]
				zero = true
			}
		} else {
			res += zhCa[site] + company[len(nStr)-i]
			zero = false
		}
	}
	return res, nil
}
