package model

import (
	"database/sql"
	"demacia/common/cachemodel"
	"fmt"
	"github.com/go-xorm/builder"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"strings"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
)

var (
	subjectFieldNames = strings.Join(cachemodel.RawFieldNames(&Subject{}, true), ",")

	cacheSubjectSubjectIdPrefix    = "cache:subject:subject:id:"
	cacheSubjectSubjectOrgIdPrefix = "cache:subject:subject:org-id:"
)

type (
	SubjectModel struct {
		*cachemodel.CachedModel
	}
	Subject struct {
		Id    int64  `db:"id"`
		Title string `db:"title"`  // 学科名
		OrgId int64  `db:"org_id"` // 机构ID
	}
	Subjects []*Subject
)

func NewSubjectModel(conn sqlx.SqlConn, cache *redis.Redis) *SubjectModel {
	return &SubjectModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"subject"."subject"`, cache),
	}
}

func (m *SubjectModel) ListSubjectByOrgId(OrgId int64) (*Subjects, error) {
	var res Subjects
	subjectListByOrgIdKey := fmt.Sprintf("%s%d", cacheSubjectSubjectOrgIdPrefix, OrgId)
	err := m.QueryRow(&res, subjectListByOrgIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		sql, args, _ := builder.Postgres().
			Select(subjectFieldNames).
			From(m.Table).
			Where(builder.Eq{"org_id": OrgId, "deleted": false}).
			ToSQL()
		return conn.QueryRows(v, sql, args...)
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (m *SubjectModel) ListSubjectByTitleAndOrgId(title string, orgId int64) (*Subjects, error) {
	query, args, _ := builder.Postgres().
		Select(subjectFieldNames).
		From(m.Table).
		Where(builder.Like{"title", title}).
		And(builder.Eq{"deleted": false, "org_id": orgId}).
		ToSQL()
	var res Subjects
	err := m.Conn.QueryRows(&res, query, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (m *SubjectModel) GetSubjectByTitleAndOrgId(title string, orgId int64) (*Subject, error) {
	query, args, _ := builder.Postgres().
		Select(subjectFieldNames).
		From(m.Table).
		Where(builder.Eq{"title": title}).
		And(builder.Eq{"deleted": false, "org_id": orgId}).
		ToSQL()
	var res Subject
	err := m.Conn.QueryRow(&res, query, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *SubjectModel) FindSubjectById(Id int64) (*Subject, error) {
	subjectIdKey := m.formatPrimary(Id)
	query, args, _ := builder.Postgres().
		Select(subjectFieldNames).
		From(m.Table).
		Where(builder.Eq{"id": Id, "deleted": false}).
		ToSQL()
	var res Subject
	err := m.Conn.QueryRow(&res, query, args...)
	if err != nil {
		//if err != cachemodel.ErrNotFound {
		//	logx.Errorf("Subject Api TeacherManage subject[model] FindSubjectById err : %s ", err.Error())
		//}
		return nil, err
	}
	_ = m.DelCache(subjectIdKey)
	return &res, nil
}

func (m *SubjectModel) Insert(Subject *Subject) (int64, error) {
	sqlString, args, _ := builder.Postgres().Insert(builder.Eq{
		"title":  Subject.Title,
		"org_id": Subject.OrgId,
	}).Into(m.Table).ToSQL()
	var LastInsertId int64
	err := m.Conn.QueryRow(&LastInsertId, sqlString+" returning id", args...)
	if err != nil {
		return 0, err
	}
	_ = m.DelCache(m.keys(Subject)...)
	return LastInsertId, nil
}

func (m *SubjectModel) Update(Subject *Subject) error {
	subjectInfo, err := m.FindSubjectById(Subject.Id)
	if err != nil || subjectInfo.OrgId != Subject.OrgId {
		return cachemodel.ErrNotFound
	}
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		sqlString, args, _ := builder.Postgres().Update(builder.Eq{
			"title": Subject.Title,
		}).Where(builder.Eq{
			"id":      Subject.Id,
			"org_id":  Subject.OrgId,
			"deleted": false,
		}).From(m.Table).ToSQL()
		return conn.Exec(sqlString, args...)
	}, m.keys(subjectInfo)...)
	return err
}
func (m *SubjectModel) Rename(Subject *Subject) error {
	subjectInfo, err := m.FindSubjectById(Subject.Id)
	if err != nil || subjectInfo.OrgId != Subject.OrgId {
		return cachemodel.ErrNotFound
	}
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		sqlString, args, _ := builder.Postgres().Update(builder.Eq{
			"title": Subject.Title,
		}).Where(builder.Eq{
			"id":      Subject.Id,
			"org_id":  Subject.OrgId,
			"deleted": false,
		}).From(m.Table).ToSQL()
		return conn.Exec(sqlString, args...)
	}, m.keys(subjectInfo)...)
	return err
}

func (m *SubjectModel) DeletedByOrgIdAndId(OrgId, Id int64) error {
	subjectInfo, err := m.FindSubjectById(Id)
	if err != nil || subjectInfo.OrgId != OrgId {
		return cachemodel.ErrNotFound
	}
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		sqlString, args, _ := builder.Postgres().Update(builder.Eq{
			"deleted": true,
		}).Where(builder.Eq{
			"id":     subjectInfo.Id,
			"org_id": subjectInfo.OrgId,
		}).From(m.Table).ToSQL()
		return conn.Exec(sqlString, args...)
	}, m.keys(subjectInfo)...)
	return err
}

func (m *SubjectModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheSubjectSubjectIdPrefix, primary)
}

func (m *SubjectModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = $1 and deleted = 0 limit 1", subjectFieldNames, m.Table)
	return conn.QueryRow(v, query, primary)
}

func (m *SubjectModel) keys(subject *Subject) []string {
	res := make([]string, 0, 2)
	if subject.OrgId != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheSubjectSubjectOrgIdPrefix, subject.OrgId))
	}
	if subject.Id != 0 {
		res = append(res, m.formatPrimary(subject.Id))
	}
	return res
}
