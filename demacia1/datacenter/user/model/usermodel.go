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
	userFieldNames          = strings.Join(cachemodel.RawFieldNames(&User{}, true), ",")
	cacheUserIdPrefix       = "cache:user:user:id:"
	cacheUserUserNamePrefix = "cache:user:user:user-name:"
)

type (
	UserModel struct {
		*cachemodel.CachedModel
	}
	User struct {
		Id       int64  `db:"id"`
		UserName string `db:"user_name"`
		Password string `db:"password"`
		Mobile   string `db:"mobile"`
		TrueName string `db:"true_name"`
	}
)

func NewUserModel(conn sqlx.SqlConn, cache *redis.Redis) *UserModel {
	return &UserModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"user"."user"`, cache),
	}
}

func (m *UserModel) FindOneById(id int64) (*User, error) {
	key := fmt.Sprintf("%s%d", cacheUserIdPrefix, id)
	var resp User
	err := m.QueryRow(&resp, key, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().From(m.Table).Select(userFieldNames).Where(builder.Eq{"id": id, "deleted": false}).Limit(1).ToSQL()
		return conn.QueryRow(v, query, args...)
	})
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

func (m *UserModel) FindOneByUserName(userName string) (*User, error) {
	key := fmt.Sprintf("%s%s", cacheUserUserNamePrefix, userName)
	var resp User
	err := m.QueryRow(&resp, key, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().From(m.Table).Select(userFieldNames).Where(builder.Eq{"user_name": userName, "deleted": false}).Limit(1).ToSQL()
		return conn.QueryRow(v, query, args...)
	})
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

func (m *UserModel) Insert(data *User) (int64, error) {
	query, args, _ := builder.Postgres().Insert(builder.Eq{
		"user_name": data.UserName,
		"mobile":    data.Mobile,
		"password":  data.Password,
		"true_name": data.TrueName,
	}).Into(m.Table).ToSQL()
	var res int64
	err := m.Conn.QueryRow(&res, query+" returning id", args...)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (m *UserModel) DeleteById(id int64) error {
	userInfo, err := m.FindOneById(id)
	if err != nil {
		return err
	}
	userNameKey := fmt.Sprintf("%s%s", cacheUserUserNamePrefix, userInfo.UserName)
	idKey := fmt.Sprintf("%s%d", cacheUserUserNamePrefix, id)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.Eq{"id": id}).ToSQL()
		return conn.Exec(query, args...)
	}, idKey, userNameKey)
	return err
}
