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
	memberFieldNames          = strings.Join(cachemodel.RawFieldNames(&Member{}, true), ",")
	cacheMemberIdPrefix       = "cache:member:member:id:"
	cacheOrgIdPrefix          = "cache:member:member:org_id:"
	cacheMemberUserNamePrefix = "cache:member:member:user-name:"
)

const (
	TeamCreator = 2
	TeamManage  = 1
	TeamMember  = 0
)

type (
	MemberModel struct {
		*cachemodel.CachedModel
	}

	Member struct {
		Id         int64  `db:"id"`          // 自增主键
		OrgId      int64  `db:"org_id"`      // 机构id
		TrueName   string `db:"true_name"`   // 昵称
		UserName   string `db:"user_name"`   // 用户名
		Password   string `db:"password"`    // 密码
		Mobile     string `db:"mobile"`      // 手机号
		Status     int8   `db:"status"`      // 状态 0：未验证手机号，1：正常，2：离职
		Deleted    bool   `db:"deleted"`     // 是否删除
		Sex        int8   `db:"sex"`         //	性别 0：女；1：男
		JoinDate   string `db:"join_date"`   // 入职日期
		LeaveDate  string `db:"leave_date"`  // 离职日期
		Avatar     string `db:"avatar"`      // 头像
		Face       string `db:"face_url"`    // 人脸url
		FaceStatus int8   `db:"face_status"` // 人脸状态
		Role       int64  `db:"role"`        // 成员身份
	}

	MemberList []*Member

	ListCondition struct {
		TrueName   string
		OrgId      int64
		Page       int
		Limit      int
		FaceStatus int8
	}
)

func NewMemberModel(conn sqlx.SqlConn, cache *redis.Redis) *MemberModel {
	return &MemberModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"member"."member"`, cache),
	}
}

// FindOneByMobile 依据手机号查询手机号是否存在
func (m *MemberModel) FindOneByMobile(mobile string, Oid int64) (*Member, error) {
	var resp Member
	toSQL, i, err := builder.Postgres().Select(memberFieldNames).From(m.Table).Where(builder.Eq{"mobile": mobile, "org_id": Oid}).And(builder.Neq{"status": 2, "deleted": true}).ToSQL()
	if err != nil {
		return nil, err
	}
	res := m.QueryRowNoCache(&resp, toSQL, i...)
	if res != nil {
		return nil, err // true
	}
	return &resp, nil // false

}

// FindOneById 依据用户id搜索用户
func (m *MemberModel) FindOneById(id int64) (*Member, error) {
	key := fmt.Sprintf("%s%d", cacheMemberIdPrefix, id)
	var resp Member
	err := m.QueryRow(&resp, key, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().From(m.Table).Select(memberFieldNames).Where(builder.Eq{"id": id, "deleted": false}).Limit(1).ToSQL()
		return conn.QueryRow(v, query, args...)
	})
	if err != nil {
		if err != cachemodel.ErrNotFound {
			logx.Errorf("get member detail error. id is %d, error is %s", id, err.Error())
		}
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

// FindOneByUserName 依据用户名搜索用户
func (m *MemberModel) FindOneByUserName(userName string) (*Member, error) {
	key := fmt.Sprintf("%s%s", cacheMemberUserNamePrefix, userName)
	var resp Member
	err := m.QueryRow(&resp, key, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().From(m.Table).Select(memberFieldNames).Where(builder.Eq{"user_name": userName, "deleted": false}).Limit(1).ToSQL()
		return conn.QueryRow(v, query, args...)
	})
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

// Insert 添加member
func (m *MemberModel) Insert(data *Member) (int64, error) {
	query, args, _ := builder.Postgres().Insert(builder.Eq{
		"org_id":    data.OrgId,
		"user_name": data.UserName,
		"mobile":    data.Mobile,
		"password":  data.Password,
		"true_name": data.TrueName,
		"join_date": time.Now().Format("2006-01-02"), // 入职时间
		"sex":       data.Sex,                        // 性别
		"avatar":    data.Avatar,                     // 头像地址
		"face_url":  data.Face,                       // 人脸地址

	}).Into(m.Table).ToSQL()
	var res int64
	err := m.Conn.QueryRow(&res, query+" returning id", args...)
	if err != nil {
		return 0, err
	}
	return res, nil
}

// SetMemberStatus 设置用户状态
func (m *MemberModel) SetMemberStatus(s int64, Mid int64) error {
	key := fmt.Sprintf("%s%d", cacheMemberIdPrefix, Mid)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"status": s}).From(m.Table).Where(builder.Eq{"id": Mid}).ToSQL()
		return conn.Exec(query, args...)
	}, key)
	return err
}

// DelMembers 删除指定的多个成员
func (m *MemberModel) DelMembers(s []int64) error {
	query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.In("id", s)).ToSQL()
	_, err := m.ExecNoCache(query, args...)
	return err
}

// Details 获取详情
func (m *MemberModel) Details(mid, orgId int64) (*Member, error) {
	// 组装sql 语句
	toSQL, i, _ := builder.Postgres().Select(memberFieldNames).From(m.Table).Where(builder.Eq{"id": mid, "org_id": orgId, "deleted": false}).ToSQL()
	var res = Member{}
	err := m.QueryRowNoCache(&res, toSQL, i...)
	if err != nil {
		return nil, err
	}
	return &res, nil

}

// FindListByOrgId 依据单位id搜索用户列表 list
func (m *MemberModel) FindListByOrgId(orgId int64) (MemberList, error) {
	key := fmt.Sprintf("%s%d", cacheOrgIdPrefix, orgId)
	var resp MemberList
	err := m.QueryRow(&resp, key, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().From(m.Table).Select(memberFieldNames).Where(builder.Eq{"org_id": orgId, "deleted": false}).OrderBy("mobile").ToSQL()
		return conn.QueryRows(v, query, args...)
	})
	return resp, err
}

// AdvancedSearchList 高级搜索列表 分页
func (m *MemberModel) AdvancedSearchList(listCondition ListCondition) (int64, MemberList, error) {
	// 过滤高级搜素条件
	where := builder.Eq{"deleted": false}.And(builder.Eq{"org_id": listCondition.OrgId})
	// 昵称模糊查询
	if listCondition.TrueName != "" {
		where = where.And(builder.Like{"true_name", listCondition.TrueName})
	}
	// 指定人脸状态 -2 查询全部
	if listCondition.FaceStatus != -2 {
		where = where.And(builder.Eq{"face_status": listCondition.FaceStatus})
	}
	// 获取总记录条数
	query, args, _ := builder.Postgres().Select("count(*)").From(m.Table).Where(where).ToSQL()
	var count int64
	err := m.QueryRowNoCache(&count, query, args...)
	if err != nil {
		return 0, nil, err
	}
	// 如果当总记录条数为0 或者 需要查询的数量 大于了总记录条数
	if count == 0 || (listCondition.Page-1)*listCondition.Limit >= int(count) {
		return count, MemberList{}, nil
	}

	// 组装sql 语句
	toSQL, i, _ := builder.Postgres().Select(memberFieldNames).From(m.Table).Where(where).OrderBy("user_name").Limit(listCondition.Limit, (listCondition.Page-1)*listCondition.Limit).ToSQL()
	// 查询
	var resp MemberList
	res := m.QueryRowsNoCache(&resp, toSQL, i...)
	if res != nil {
		return 0, nil, res
	}
	return count, resp, nil
}

// DeleteById 删除成员
func (m *MemberModel) DeleteById(id int64) error {
	memberInfo, err := m.FindOneById(id)
	if err != nil {
		return err
	}
	userNameKey := fmt.Sprintf("%s%s", cacheMemberUserNamePrefix, memberInfo.UserName)
	idKey := fmt.Sprintf("%s%d", cacheMemberUserNamePrefix, id)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.Eq{"id": id}).ToSQL()
		return conn.Exec(query, args...)
	}, idKey, userNameKey)
	return err
}

// DeleteByOrgId 删除单位下所有成员
func (m *MemberModel) DeleteByOrgId(orgId int64) error {
	query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.Eq{"org_id": orgId, "deleted": false}).ToSQL()
	_, err := m.ExecNoCache(query, args...)
	return err
}

// UpdateMemberInfo 修改成员详情
func (m *MemberModel) UpdateMemberInfo(mid int64, data Member) error {
	set := builder.Eq{}

	if data.TrueName != "" {
		set["true_name"] = data.TrueName
	}
	if data.UserName != "" {
		set["user_name"] = data.UserName
	}
	if data.Mobile != "" {
		set["mobile"] = data.Mobile
	}
	if data.Avatar != "" {
		set["avatar"] = data.Avatar
	}
	if data.Sex > -1 {
		set["sex"] = data.Sex
	}
	if data.Face != "" {
		set["face_url"] = data.Face
	}
	if len(set) < 0 {
		return nil
	}

	key := fmt.Sprintf("%s%d", cacheMemberIdPrefix, mid)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(set).From(m.Table).Where(builder.Eq{"id": mid}).ToSQL()
		return conn.Exec(query, args...)
	}, key)
	return err
}

// SetMemberRole 设置成员角色
func (m *MemberModel) SetMemberRole(mid int64, orgId int64, role int64) error {
	key := fmt.Sprintf("%s%d", cacheOrgIdPrefix, orgId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"role": role}).From(m.Table).Where(builder.Eq{"id": mid, "org_id": orgId}).ToSQL()
		return conn.Exec(query, args...)
	}, key)
	return err
}

// LikeByTrueName 依据姓名匹配列表
func (m MemberModel) LikeByTrueName(truename string, orgId int64) (MemberList, error) {

	toSQL, i, err := builder.Postgres().Select(memberFieldNames).From(m.Table).Where(builder.Eq{"org_id": orgId, "deleted": false}).And(builder.Like{"true_name", truename}).ToSQL()
	if err != nil {
		return nil, err
	}
	var resp MemberList
	err = m.QueryRowsNoCache(&resp, toSQL, i...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
