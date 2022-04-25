package model

import (
	"database/sql"
	"fmt"
	"strings"

	"demacia/common/basefunc"
	"demacia/common/cachemodel"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"xorm.io/builder"
)

var (
	departmentMemberFieldNames       = strings.Join(cachemodel.RawFieldNames(&DepartmentMember{}, true), ",")
	cacheDepartmentMemberOrgIdPrefix = "cache:department:department_member:org-id:"
)

type (
	DepartmentMemberModel struct {
		*cachemodel.CachedModel
	}
	DepartmentMember struct {
		Id           int64  `db:"id"`            // 自增主键
		MemberId     int64  `db:"member_id"`     // 部门用户id
		DepartmentId int64  `db:"department_id"` // 部门id
		OrgId        int64  `db:"org_id"`        // 机构ID
		TrueName     string `db:"true_name"`     // 人员姓名
		Mobile       string `db:"mobile"`        // 手机号
	}
	DepartmentMembers []*DepartmentMember
)

func NewDepartmentMemberModel(conn sqlx.SqlConn, cache *redis.Redis) *DepartmentMemberModel {
	return &DepartmentMemberModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"department"."department_member"`, cache),
	}
}

func (m *DepartmentMemberModel) GetDepartmentMembersByOrgId(orgId int64) (DepartmentMembers, error) {
	var res DepartmentMembers
	err := m.QueryRow(&res, m.formatPrimary(orgId), func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().Select(departmentMemberFieldNames).From(m.Table).Where(builder.Eq{"org_id": orgId, "deleted": false}).ToSQL()
		return conn.QueryRows(v, query, args...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *DepartmentMemberModel) GetDepartmentIdsByOrgIdAndMemberId(orgId, memberId int64) ([]int64, error) {
	var res []int64
	query, args, _ := builder.Postgres().Select("department_id").From(m.Table).Where(builder.Eq{"org_id": orgId, "member_id": memberId, "deleted": false}).ToSQL()
	err := m.Conn.QueryRows(&res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *DepartmentMemberModel) GetMembersByOrgIdAndDepartmentId(orgId, departmentId int64) (DepartmentMembers, error) {
	var res DepartmentMembers
	query, args, _ := builder.Postgres().Select(departmentMemberFieldNames).From(m.Table).Where(builder.Eq{"org_id": orgId, "department_id": departmentId, "deleted": false}).ToSQL()
	err := m.Conn.QueryRows(&res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *DepartmentMemberModel) DeleteByOrgIdAndDepartmentIdAndMemberIds(orgId, departmentId int64, memberIds []int64) error {
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).Where(builder.Eq{"org_id": orgId, "department_id": departmentId, "deleted": false}).And(builder.In("member_id", memberIds)).From(m.Table).ToSQL()
		return conn.Exec(query, args...)
	}, m.formatPrimary(orgId))
	return err
}

func (m *DepartmentMemberModel) DeleteByByOrgIdAndMemberId(orgId, memberId int64) error {
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).Where(builder.Eq{"org_id": orgId, "member_id": memberId, "deleted": false}).From(m.Table).ToSQL()
		return conn.Exec(query, args...)
	}, m.formatPrimary(orgId))
	return err
}

func (m *DepartmentMemberModel) UpdateMemberInfo(data *DepartmentMember) error {
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"org_id": data.OrgId, "true_name": data.TrueName, "mobile": data.Mobile}).Where(builder.Eq{"member_id": data.MemberId, "deleted": false}).From(m.Table).ToSQL()
		return conn.Exec(query, args...)
	}, m.formatPrimary(data.OrgId))
	return err
}

func (m *DepartmentMemberModel) BatchInsert(data DepartmentMembers) error {
	if len(data) == 0 {
		return nil
	}
	query, args, err := basefunc.BatchInsertString(m.Table, data)
	if err != nil {
		return err
	}
	_, err = m.ExecNoCache(query, args...)
	return err
}

func (m *DepartmentMemberModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheDepartmentMemberOrgIdPrefix, primary)
}
