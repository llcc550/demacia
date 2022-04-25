package model

import (
	"strings"

	"demacia/common/cachemodel"

	"github.com/go-xorm/builder"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
)

type (
	SmsTraffic struct {
		Id          int64  `db:"id" json:"id,omitempty"`
		Mobile      string `db:"mobile" json:"mobile,omitempty"`            // 手机号
		ContentType string `db:"content_type" json:"contentType,omitempty"` // 短信类型
		Content     string `db:"content" json:"content,omitempty"`          // 短信内容
		Provider    string `db:"provider" json:"provider,omitempty"`        // 运营商
		Error       string `db:"error"  json:"error,omitempty"`             // 错误
	}
	SmsTrafficList []*SmsTraffic

	SmsTrafficModel struct {
		*cachemodel.CachedModel
	}
)

func NewSmsTrafficModel(conn sqlx.SqlConn, rds *redis.Redis) *SmsTrafficModel {
	return &SmsTrafficModel{CachedModel: cachemodel.NewCachedModel(conn, `"sms"."sms_traffic"`, rds)}
}

func (st *SmsTrafficModel) Insert(message *SmsTraffic) error {
	sql, args, _ := builder.Postgres().Insert(builder.Eq{
		"mobile":       message.Mobile,
		"content_type": message.ContentType,
		"content":      message.Content,
		"provider":     message.Provider,
		"error":        message.Error,
	}).Into(st.Table).ToSQL()
	sql = strings.Replace(sql, `\`, "", -1)
	_, err := st.Conn.Exec(sql, args...)
	return err
}

func (st *SmsTrafficModel) InsertBatch(data *SmsTrafficList) error {
	if len(*data) == 0 {
		return nil
	}
	//TODO 批量插入
	//tableName := st.table + "_" + time.Now().Format("200601")
	//st.createTable(tableName)
	//values, args, err := utils.BatchInsertString(*data)
	//if err != nil {
	//	return err
	//}
	//_, err = st.conn.Exec(`INSERT INTO `+st.table+` (`+dormQueryRows+`) VALUES `+values, args...)
	return nil
}
