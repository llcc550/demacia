package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"demacia/common/cachemodel"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"xorm.io/builder"
)

var (
	organizationFieldNames    = strings.Join(cachemodel.RawFieldNames(&Organization{}, true), ",")
	cacheOrganizationIdPrefix = "cache:organization:organization:id:"
)

type (
	OrganizationModel struct {
		*cachemodel.CachedModel
	}
	Organization struct {
		Id                    int64  `db:"id"`
		CreateTime            int64  `db:"create_time"`              // 创建时间
		ProvinceId            int64  `db:"province_id"`              // 省
		CityId                int64  `db:"city_id"`                  // 市
		AreaId                int64  `db:"area_id"`                  // 区
		AgentId               int64  `db:"agent_id"`                 // 代理商ID，值为0代表无代理商
		ManagerMemberId       int64  `db:"manager_member_id"`        // 拥有者ID
		TermId                int64  `db:"term_id"`                  // 当前学期
		Title                 string `db:"title"`                    // 机构名称
		ActivateDate          string `db:"activate_date"`            // 启用日期
		ExpireDate            string `db:"expire_date"`              // 过期日期
		Addr                  string `db:"addr"`                     // 单位地址
		Msg                   string `db:"msg"`                      // 备注
		TrueName              string `db:"true_name"`                // 联系人
		Mobile                string `db:"mobile"`                   // 联系方式
		AreaTitle             string `db:"area_title"`               // 省-市-区。冗余数据
		ManagerMemberUserName string `db:"manager_member_user_name"` // 拥有者账号
		OrgType               int8   `db:"org_type"`                 // 1：学校，其他值待定
		OrgStatus             int8   `db:"org_status"`               // 0：未启用，1：正常使用，-1：已过期，-2：禁用
	}
	Organizations []*Organization
	ListCond      struct {
		ProvinceId int64
		CityId     int64
		AreaId     int64
		Title      string
	}
)

func NewOrganizationModel(conn sqlx.SqlConn, cache *redis.Redis) *OrganizationModel {
	return &OrganizationModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"organization"."organization"`, cache),
	}
}

func (m *OrganizationModel) List(page, limit int, cond *ListCond) (int64, Organizations, error) {
	where := builder.Eq{"deleted": false}.And()
	if cond.Title != "" {
		where = where.And(builder.Like{"title", cond.Title})
	}
	if cond.ProvinceId != 0 {
		where = where.And(builder.Eq{"province_id": cond.ProvinceId})
		if cond.CityId != 0 {
			where = where.And(builder.Eq{"city_id": cond.CityId})
			if cond.AreaId != 0 {
				where = where.And(builder.Eq{"area_id": cond.AreaId})
			}
		}
	}

	query, args, _ := builder.Postgres().Select("count(*)").From(m.Table).Where(where).ToSQL()
	var count int64
	err := m.Conn.QueryRow(&count, query, args...)
	if err != nil {
		return 0, nil, err
	}
	if count == 0 || (page-1)*limit >= int(count) {
		return count, Organizations{}, nil
	}
	query, args, _ = builder.Postgres().Select(organizationFieldNames).From(m.Table).Where(where).OrderBy("id").Limit(limit, (page-1)*limit).ToSQL()
	var list Organizations
	err = m.Conn.QueryRows(&list, query, args...)
	if err != nil {
		return 0, nil, err
	}
	return count, list, nil
}

func (m *OrganizationModel) FindOneById(orgId int64) (*Organization, error) {
	var resp Organization
	err := m.QueryRow(&resp, m.formatPrimary(orgId), func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().Select(organizationFieldNames).From(m.Table).Where(builder.Eq{"id": orgId, "deleted": false}).Limit(1).ToSQL()
		return conn.QueryRow(v, query, args...)
	})
	if err != nil {
		if err != cachemodel.ErrNotFound {
			logx.Errorf("get organization detail error. id is %d, error is %s", orgId, err.Error())
		}
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

func (m *OrganizationModel) FindIdByTitle(title string) (int64, error) {
	var resp int64
	query, args, _ := builder.Postgres().Select(organizationFieldNames).From(m.Table).Where(builder.Eq{"title": title, "deleted": false}).Limit(1).ToSQL()
	err := m.Conn.QueryRow(&resp, query, args...)
	if err != nil {
		return 0, cachemodel.ErrNotFound
	}
	return resp, nil
}

func (m *OrganizationModel) UpdateManagerMember(orgId, managerMemberId int64, managerMemberUserName string) error {
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"manager_member_id": managerMemberId, "manager_member_user_name": managerMemberUserName}).From(m.Table).Where(builder.Eq{"id": orgId, "deleted": false}).ToSQL()
		return conn.Exec(query, args...)
	}, m.formatPrimary(orgId))
	return err
}

func (m *OrganizationModel) DeleteById(orgId int64) error {
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.Eq{"id": orgId}).ToSQL()
		return conn.Exec(query, args...)
	}, m.formatPrimary(orgId))
	return err
}

func (m *OrganizationModel) DeleteByIdForce(orgId int64) error {
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Delete(builder.Eq{"deleted": false}).From(m.Table).Where(builder.Eq{"id": orgId}).ToSQL()
		return conn.Exec(query, args...)
	}, m.formatPrimary(orgId))
	return err
}

func (m *OrganizationModel) Insert(data *Organization) (int64, error) {
	orgStatus := 1
	timeNow := time.Now().Format("2006-01-02")
	if data.ActivateDate > timeNow {
		orgStatus = 0
	} else if data.ExpireDate < timeNow {
		orgStatus = -1
	}
	query, args, _ := builder.Postgres().Insert(builder.Eq{
		"title":         data.Title,
		"org_type":      data.OrgType,
		"area_id":       data.AreaId,
		"org_status":    orgStatus,
		"activate_date": data.ActivateDate,
		"expire_date":   data.ExpireDate,
		"create_time":   time.Now().Unix(),
		"province_id":   data.ProvinceId,
		"city_id":       data.CityId,
		"area_title":    data.AreaTitle,
		"agent_id":      data.AgentId,
		"addr":          data.Addr,
		"msg":           data.Msg,
		"term_id":       data.TermId,
		"true_name":     data.TrueName,
		"mobile":        data.Mobile,
	}).Into(m.Table).ToSQL()
	var res int64
	err := m.Conn.QueryRow(&res, query+" returning id", args...)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (m *OrganizationModel) Update(data *Organization) error {
	orgStatus := 1
	timeNow := time.Now().Format("2006-01-02")
	if data.ActivateDate > timeNow {
		orgStatus = 0
	} else if data.ExpireDate < timeNow {
		orgStatus = -1
	}
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{
			"title":         data.Title,
			"area_id":       data.AreaId,
			"org_status":    orgStatus,
			"activate_date": data.ActivateDate,
			"expire_date":   data.ExpireDate,
			"province_id":   data.ProvinceId,
			"city_id":       data.CityId,
			"area_title":    data.AreaTitle,
			"addr":          data.Addr,
			"msg":           data.Msg,
			"true_name":     data.TrueName,
			"mobile":        data.Mobile,
		}).From(m.Table).Where(builder.Eq{"id": data.Id, "deleted": false}).ToSQL()
		return conn.Exec(query, args...)
	}, m.formatPrimary(data.Id))
	return err
}

func (m *OrganizationModel) formatPrimary(orgId int64) string {
	return fmt.Sprintf("%s%d", cacheOrganizationIdPrefix, orgId)
}
