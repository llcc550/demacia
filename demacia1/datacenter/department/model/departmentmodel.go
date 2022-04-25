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
	departmentFieldNames       = strings.Join(cachemodel.RawFieldNames(&Department{}, true), ",")
	cacheDepartmentIdPrefix    = "cache:department:department:id:"
	cacheDepartmentOrgIdPrefix = "cache:department:department:org-id:"
	cacheDepartmentTitlePrefix = "cache:department:department:title:"
)

type (
	DepartmentModel struct {
		*cachemodel.CachedModel
	}
	Department struct {
		Id          int64  `db:"id"`           // 自增主键
		OrgId       int64  `db:"org_id"`       // 单位id
		Title       string `db:"title"`        // 部门名称
		Sort        int64  `db:"sort"`         // 自定义排序
		MemberCount int64  `db:"member_count"` // 人员总数
	}
	Departments []*Department
)

func NewDepartmentModel(conn sqlx.SqlConn, cache *redis.Redis) *DepartmentModel {
	return &DepartmentModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"department"."department"`, cache),
	}
}

func (m *DepartmentModel) GetDepartmentsByOrgId(orgId int64) (Departments, error) {
	var res Departments
	key := fmt.Sprintf("%s%d", cacheDepartmentOrgIdPrefix, orgId)
	fmt.Println(key)
	err := m.QueryRow(&res, key, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().
			Select(departmentFieldNames).
			From(m.Table).
			Where(builder.Eq{"org_id": orgId, "deleted": false}).
			OrderBy("sort ASC, id DESC").ToSQL()
		return conn.QueryRows(v, query, args...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *DepartmentModel) GetDepartmentsByOrgIdWithTitle(orgId int64, title string) (Departments, error) {
	if title == "" {
		return m.GetDepartmentsByOrgId(orgId)
	}
	var res Departments
	query, args, _ := builder.Postgres().
		Select(departmentFieldNames).
		From(m.Table).
		Where(builder.Eq{"org_id": orgId, "deleted": false}).And(builder.Like{"title", title}).
		OrderBy("sort ASC, id DESC").ToSQL()
	err := m.Conn.QueryRows(&res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *DepartmentModel) GetDepartmentById(Id int64) (*Department, error) {
	var res Department
	key := m.formatPrimary(Id)
	err := m.QueryRow(&res, key, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().Select(departmentFieldNames).From(m.Table).Where(builder.Eq{"id": Id, "deleted": false}).ToSQL()
		return conn.QueryRow(v, query, args...)
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *DepartmentModel) GetDepartmentByTitle(orgId int64, title string) (*Department, error) {
	var resp Department
	key := fmt.Sprintf("%s%d-%s", cacheDepartmentTitlePrefix, orgId, title)
	err := m.QueryRowIndex(&resp, key, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query, args, _ := builder.Postgres().From(m.Table).Select(departmentFieldNames).Where(builder.Eq{"org_id": orgId, "title": title, "deleted": false}).Limit(1).ToSQL()
		if err := conn.QueryRow(&resp, query, args...); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

func (m *DepartmentModel) InsertDepartment(data *Department) (int64, error) {
	query, args, _ := builder.Postgres().Insert(builder.Eq{
		"org_id":       data.OrgId,
		"title":        data.Title,
		"sort":         data.Sort,
		"member_count": data.MemberCount,
	}).Into(m.Table).ToSQL()

	var id int64
	err := m.Conn.QueryRow(&id, query+" returning id", args...)
	if err != nil {
		return 0, err
	}
	_ = m.initDepartmentSort(data.OrgId)
	_ = m.DelCache(m.keys(data)...)
	return id, nil
}

func (m *DepartmentModel) DeleteById(id int64) error {
	department, err := m.GetDepartmentById(id)
	if err != nil {
		return cachemodel.ErrNotFound
	}
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).Where(builder.Eq{"id": id}).From(m.Table).ToSQL()
		return conn.Exec(query, args...)
	}, m.keys(department)...)
	return m.initDepartmentSort(department.OrgId)
}

func (m *DepartmentModel) UpdateDepartmentTitle(id int64, title string) error {
	department, err := m.GetDepartmentById(id)
	if err != nil {
		return cachemodel.ErrNotFound
	}
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"title": title}).Where(builder.Eq{"id": id, "deleted": false}).From(m.Table).ToSQL()
		return conn.Exec(query, args...)
	}, m.keys(department)...)
	return err
}

func (m *DepartmentModel) UpdateDepartmentMemberCount(id, memberCount int64) error {
	department, err := m.GetDepartmentById(id)
	if err != nil {
		return cachemodel.ErrNotFound
	}
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"member_count": memberCount}).Where(builder.Eq{"id": id, "deleted": false}).From(m.Table).ToSQL()
		return conn.Exec(query, args...)
	}, m.keys(department)...)
	return err
}

func (m *DepartmentModel) UpdateDepartmentSort(orgId int64, data Departments) error {
	valueSql := ""
	for _, item := range data {
		valueSql += fmt.Sprintf(` ( %d , %d ),`, item.Id, item.Sort)
	}
	valueSql = strings.Trim(valueSql, ",")
	sortSql := fmt.Sprintf(`UPDATE %s department SET sort = tmp.sort FROM ( VALUES %s ) AS tmp ( id, sort ) WHERE department.id = tmp.id`, m.Table, valueSql)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.Exec(sortSql)
	}, fmt.Sprintf("%s%d", cacheDepartmentOrgIdPrefix, orgId))
	if err != nil {
		return err
	}
	return m.initDepartmentSort(orgId)
}

func (m *DepartmentModel) initDepartmentSort(orgId int64) error {
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf(`UPDATE %s t1 SET sort = t2.newSort FROM ( SELECT id, ROW_NUMBER () OVER ( ORDER BY sort ASC,id DESC ) AS newSort FROM %s WHERE org_id = $1 AND deleted = FALSE ) t2 WHERE t1.id = t2.id`, m.Table, m.Table)
		return conn.Exec(query, orgId)
	}, fmt.Sprintf("%s%d", cacheDepartmentOrgIdPrefix, orgId))
	return err
}

func (m *DepartmentModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheDepartmentIdPrefix, primary)
}

func (m *DepartmentModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = $1 and deleted = 0 limit 1", departmentFieldNames, m.Table)
	return conn.QueryRow(v, query, primary)
}

func (m *DepartmentModel) keys(department *Department) []string {
	res := make([]string, 0, 3)
	if department.OrgId != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheDepartmentOrgIdPrefix, department.OrgId))
	}
	if department.OrgId != 0 && department.Title != "" {
		res = append(res, fmt.Sprintf("%s%d:%s", cacheDepartmentTitlePrefix, department.OrgId, department.Title))
	}
	if department.Id != 0 {
		res = append(res, m.formatPrimary(department.Id))
	}
	return res
}
