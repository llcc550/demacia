package model

import (
	"database/sql"
	"fmt"
	"strings"

	"demacia/common/cachemodel"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"xorm.io/builder"
)

var (
	studentFieldNames                 = strings.Join(cachemodel.RawFieldNames(&Student{}, true), ",")
	cacheStudentListByClassIdPrefix   = "cache:student:student:class-id:"
	cacheStudentListByOrgIdPrefix     = "cache:student:student:org-id:"
	cacheStudentInfoByStudentIdPrefix = "cache:student:student:id:"
)

type (
	StudentModel struct {
		*cachemodel.CachedModel
	}
	Student struct {
		Id            int64  `db:"id"`
		OrgId         int64  `db:"org_id"`
		ClassId       int64  `db:"class_id"`
		TrueName      string `db:"true_name"`
		UserName      string `db:"user_name"`
		Deleted       bool   `db:"deleted"`
		StageId       int64  `db:"stage_id"`
		GradeId       int64  `db:"grade_id"`
		ClassFullName string `db:"class_full_name"`
		Password      string `db:"password"`
		Face          string `db:"face"`
		Avatar        string `db:"avatar"`
		FaceStatus    int64  `db:"face_status"`
		IdNumber      string `db:"id_number"`
		CardNumber    string `db:"card_number"`
		Sex           int8   `json:"sex"`
	}
	ListCondition struct {
		OrgId       int64
		StageId     int64
		GradeId     int64
		ClassId     int64
		StudentName string
		Page        int
		Limit       int
		FaceStatus  int8
	}
)

func NewStudentModel(conn sqlx.SqlConn, cache *redis.Redis) *StudentModel {
	return &StudentModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"student"."student"`, cache),
	}
}

func (m *StudentModel) FindOneByOrgIdAndUserName(orgId int64, userName string) (*Student, error) {
	var resp Student
	query, args, _ := builder.Postgres().Select(studentFieldNames).From(m.Table).Where(builder.Eq{"org_id": orgId, "user_name": userName, "deleted": false}).Limit(1).ToSQL()
	err := m.Conn.QueryRow(&resp, query, args...)
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

func (m *StudentModel) InsertOne(data *Student) (int64, error) {
	keys := m.keys(data)
	query, args, _ := builder.Postgres().Insert(builder.Eq{
		"org_id":          data.OrgId,
		"stage_id":        data.StageId,
		"grade_id":        data.GradeId,
		"class_id":        data.ClassId,
		"class_full_name": data.ClassFullName,
		"user_name":       data.UserName,
		"true_name":       data.TrueName,
		"password":        data.Password,
		"avatar":          data.Avatar,
		"face":            data.Face,
		"face_status":     data.FaceStatus,
		"id_number":       data.IdNumber,
		"card_number":     data.CardNumber,
		"sex":             data.Sex,
	}).Into(m.Table).ToSQL()
	var res int64
	err := m.Conn.QueryRow(&res, query+" returning id", args...)
	if err != nil {
		return 0, err
	}
	_ = m.DelCache(keys...)
	return res, nil
}

func (m *StudentModel) UpdateOne(data *Student) error {
	keys := m.keys(data)
	query, args, _ := builder.Postgres().Update(builder.Eq{
		"org_id":          data.OrgId,
		"stage_id":        data.StageId,
		"grade_id":        data.GradeId,
		"class_id":        data.ClassId,
		"class_full_name": data.ClassFullName,
		"user_name":       data.UserName,
		"true_name":       data.TrueName,
		"password":        data.Password,
		"avatar":          data.Avatar,
		"face":            data.Face,
		"face_status":     data.FaceStatus,
		"id_number":       data.IdNumber,
		"card_number":     data.CardNumber,
		"sex":             data.Sex,
	}).From(m.Table).Where(builder.Eq{"id": data.Id}).ToSQL()
	_, err := m.ExecNoCache(query, args...)
	_ = m.DelCache(keys...)
	return err
}

func (m *StudentModel) DeleteByClassId(classId int64) error {
	studentListByClassIdKey := fmt.Sprintf("%s%d", cacheStudentListByClassIdPrefix, classId)
	query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.Eq{"class_id": classId, "deleted": false}).ToSQL()
	_, err := m.ExecNoCache(query, args...)
	_ = m.DelCache(studentListByClassIdKey)
	return err
}

func (m *StudentModel) DeleteByStudentIds(studentIds []int64) error {
	_, err := m.Exec(func(conn sqlx.SqlConn) (sql.Result, error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.In("id", studentIds)).ToSQL()
		return conn.Exec(query, args...)
	}, m.idKeys(studentIds)...)
	return err
}

func (m *StudentModel) FindOneByIdNumber(idNumber string) (*Student, error) {
	var resp Student
	query, args, _ := builder.Postgres().Select(studentFieldNames).From(m.Table).Where(builder.Eq{"id_number": idNumber, "deleted": false}).Limit(1).ToSQL()
	err := m.Conn.QueryRow(&resp, query, args...)
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

func (m *StudentModel) FindOneById(id int64) (*Student, error) {
	var resp Student
	studentInfoByStudentIdKey := m.formatPrimary(id)
	err := m.QueryRow(&resp, studentInfoByStudentIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().Select(studentFieldNames).From(m.Table).Where(builder.Eq{"id": id, "deleted": false}).Limit(1).ToSQL()
		return conn.QueryRow(v, query, args...)
	})
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

func (m *StudentModel) FindListByConditions(req *ListCondition) ([]Student, int, error) {
	var res []Student
	var count int
	eq := builder.And()
	if req.OrgId != -1 {
		eq = eq.And(builder.Eq{"org_id": req.OrgId})
	}
	if req.StageId != -1 {
		eq = eq.And(builder.Eq{"stage_id": req.StageId})
	}
	if req.GradeId != -1 {
		eq = eq.And(builder.Eq{"grade_id": req.GradeId})
	}
	if req.ClassId != -1 {
		eq = eq.And(builder.Eq{"class_id": req.ClassId})
	}
	if req.StudentName != "" {
		eq = eq.And(builder.Like{"true_name", req.StudentName})
	}
	if req.FaceStatus != -1 {
		eq = eq.And(builder.Eq{"face_status": req.FaceStatus})
	}
	query := builder.Postgres().Select(studentFieldNames).From(m.Table).Where(eq).And(builder.Eq{"deleted": false})
	if req.Page > 0 || req.Limit > 0 {
		query = query.Limit(req.Limit, (req.Page-1)*req.Limit).OrderBy("id")
	}
	sqlQ, argsQ, _ := query.ToSQL()
	sqlC, argsC, _ := builder.Postgres().Select("COUNT(*)").From(m.Table).Where(eq).And(builder.Eq{"deleted": false}).ToSQL()
	err := m.Conn.QueryRows(&res, sqlQ, argsQ...)
	err = m.Conn.QueryRow(&count, sqlC, argsC...)
	return res, count, err
}

func (m *StudentModel) FindListByClassId(classId int64) ([]Student, error) {
	var res []Student
	studentListByClassIdKey := fmt.Sprintf("%s%d", cacheStudentListByClassIdPrefix, classId)
	err := m.QueryRow(&res, studentListByClassIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		sqlQ, argsQ, _ := builder.Postgres().Select(studentFieldNames).From(m.Table).Where(builder.Eq{"class_id": classId}).And(builder.Eq{"deleted": false}).ToSQL()
		return conn.QueryRows(v, sqlQ, argsQ...)
	})
	return res, err
}

func (m *StudentModel) GetClassStudentCount(classId int64) (int64, error) {
	var count int64
	sqlQ, argsQ, _ := builder.Postgres().Select("COUNT ( * )").From(m.Table).Where(builder.Eq{"class_id": classId}).And(builder.Eq{"deleted": false}).ToSQL()
	err := m.Conn.QueryRow(&count, sqlQ, argsQ...)
	return count, err
}

func (m *StudentModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheStudentInfoByStudentIdPrefix, primary)
}

func (m *StudentModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = $1 and deleted = 0 limit 1", studentFieldNames, m.Table)
	return conn.QueryRow(v, query, primary)
}

func (m *StudentModel) keys(data *Student) []string {
	res := make([]string, 0, 4)
	if data.OrgId != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheStudentListByOrgIdPrefix, data.OrgId))
	}
	if data.ClassId != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheStudentListByClassIdPrefix, data.ClassId))
	}
	if data.Id != 0 {
		res = append(res, m.formatPrimary(data.Id))
	}
	return res
}

func (m *StudentModel) idKeys(ids []int64) []string {
	res := make([]string, 0, len(ids))
	for _, id := range ids {
		res = append(res, m.formatPrimary(id))
	}
	return res
}
