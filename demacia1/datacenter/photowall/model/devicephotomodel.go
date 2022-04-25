package model

import (
	"database/sql"
	"demacia/common/cachemodel"
	"demacia/datacenter/photowall/common"
	"fmt"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"strings"
	"time"
	"xorm.io/builder"
)

var (
	devicePhotoFieldNames = strings.Join(cachemodel.RawFieldNames(&DevicePhoto{}, true), ",")

	cachePhotowallDevicePhotoIdPrefix = "cache:photowall:devicePhoto:id:"
)

type (
	DevicePhotoModel struct {
		*cachemodel.CachedModel
	}

	DevicePhoto struct {
		Id                   int64          `db:"id"`
		OrgId                int64          `db:"org_id"`                 // 机构id
		DeviceId             int64          `db:"device_id"`              // 设备id
		PhotoId              int64          `db:"photo_id"`               // 资源id
		PublishStartTime     sql.NullString `db:"publish_start_time"`     // 发布开始时间
		PublishEndTime       sql.NullString `db:"publish_end_time"`       // 发布结束时间
		ScreensaverStartTime sql.NullString `db:"screensaver_start_time"` // 屏保开始时间
		ScreensaverEndTime   sql.NullString `db:"screensaver_end_time"`   // 屏保结束时间
		LockScreenStartTime  sql.NullString `db:"lock_screen_start_time"` // 锁屏开始时间
		LockScreenEndTime    sql.NullString `db:"lock_screen_end_time"`   // 锁屏结束时间
		ToppingStartTime     sql.NullString `db:"topping_start_time"`     // 置顶开始时间
		ToppingEndTime       sql.NullString `db:"topping_end_time"`       // 置顶结束时间
		ScreenSaverWaitTime  int            `db:"screensaver_wait_time"`  // 屏保等待时间
		PhotoUrl             string         `db:"photo_url"`              // 资源地址
		DevicePhotoTitle     string         `db:"device_photo_title"`     // 资源名称
		UpdatedTime          string         `db:"updated_time"`           // 更新时间
	}
	DevicePhotos []DevicePhoto
	PhotoTime    struct {
		ToppingEndTime     sql.NullString `db:"topping_end_time"`
		PublishEndTime     sql.NullString `db:"publish_end_time"`
		ScreensaverEndTime sql.NullString `db:"screensaver_end_time"`
		LockScreenEndTime  sql.NullString `db:"lock_screen_end_time"`
	}
	ListByIdReq struct {
		Id    int64
		Title string // 资源名称
		OrgId int64  // 机构id
		Page  int
		Limit int
	}
)

func NewDevicePhotoModel(conn sqlx.SqlConn, cache *redis.Redis) *DevicePhotoModel {
	return &DevicePhotoModel{cachemodel.NewCachedModel(conn, `"photowall"."device_photo"`, cache)}
}

func (m *DevicePhotoModel) Insert(data *DevicePhoto) (int64, error) {
	toSQL, i, err := builder.Postgres().Insert(builder.Eq{
		"org_id":                 data.OrgId,
		"device_id":              data.DeviceId,
		"photo_id":               data.PhotoId,
		"device_photo_title":     data.DevicePhotoTitle,
		"screensaver_start_time": data.LockScreenStartTime,
		"screensaver_end_time":   data.ScreensaverEndTime,
		"lock_screen_start_time": data.LockScreenStartTime,
		"lock_screen_end_time":   data.LockScreenEndTime,
		"topping_start_time":     data.ToppingStartTime,
		"topping_end_time":       data.ToppingEndTime,
		"photo_url":              data.PhotoUrl,
		"publish_start_time":     data.PublishStartTime,
		"publish_end_time":       data.PublishEndTime,
	}).Into(m.Table).ToSQL()
	if err != nil {
		logx.Errorf("devicephoto insert err:%s", err.Error())
		return 0, err
	}
	var LastInsertId int64
	err = m.Conn.QueryRow(&LastInsertId, toSQL+" returning id", i...)
	if err != nil {
		logx.Errorf("devicephoto insert err:%s", err.Error())
		return 0, err
	}
	return LastInsertId, nil
}

// ListByDeviceId 根据设备Id，分页获取资源列表
func (m *DevicePhotoModel) ListByDeviceId(data *ListByIdReq) (DevicePhotos, int, error) {
	var res DevicePhotos
	var total int
	var offset int
	if data.Limit > 0 && data.Limit < 100 {
		offset = (data.Page - 1) * data.Limit
	}
	eq := builder.And(builder.Eq{"deleted": false, "org_id": data.OrgId, "device_id": data.Id})
	if data.Title != "" {
		eq = eq.And(builder.Like{"device_photo_title", data.Title})
	}
	query, args, _ := builder.Postgres().
		Select(devicePhotoFieldNames).
		From(m.Table).
		Where(eq).
		Limit(data.Limit, offset).
		OrderBy("id DESC ").
		ToSQL()
	err := m.QueryRowsNoCache(&res, query, args...)
	if err != nil {
		logx.Errorf("devicephoto List sql :%s err:%s", query, err.Error())
		return nil, 0, err
	}
	totalQuery, totalArgs, _ := builder.Postgres().
		Select("COUNT(id)").
		From(m.Table).
		Where(eq).
		ToSQL()
	err = m.QueryRowNoCache(&total, totalQuery, totalArgs...)
	if err != nil {
		return nil, 0, err
	}
	return res, total, nil
}

func (m *DevicePhotoModel) FindPhotoExpirationTimeByOrgIdAndPhotoId(orgId, photoId int64) (*PhotoTime, error) {

	return nil, nil
}

// FindPhotoByPhotoIdAndDeviceId 根据设备Id&资源Id获取配置详情
func (m *DevicePhotoModel) FindPhotoByPhotoIdAndDeviceId(data *DevicePhoto) (*DevicePhoto, error) {
	//eq := builder.And(builder.Eq{"deleted": false})
	var res DevicePhoto
	query, args, _ := builder.Postgres().
		Select(devicePhotoFieldNames).
		From(m.Table).
		Where(builder.Eq{
			"photo_id":  data.PhotoId,
			"device_id": data.DeviceId,
			"deleted":   false,
		}).
		OrderBy("id DESC ").
		ToSQL()
	err := m.QueryRowNoCache(&res, query, args...)
	if err != nil {
		logx.Errorf("devicephoto FindPhotoByPhotoId sql :%s err:%s", query, err.Error())
		return nil, err
	}
	return &res, nil
}

func (m *DevicePhotoModel) Update(data *DevicePhoto) error {
	photowallPhotosIdKey := fmt.Sprintf("%s%v", cachePhotowallPhotosIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where id = $1", m.Table, photosFieldNames)
		return conn.Exec(query, data.Id, data.DevicePhotoTitle, data.PhotoUrl, data.OrgId, data.UpdatedTime)
	}, photowallPhotosIdKey)
	return err
}
func (m *DevicePhotoModel) UpdateTimeById(data *DevicePhoto) error {
	eq := builder.And(builder.Eq{
		"updated_time": time.Now().Format("2006-01-02 15:04:05"),
	})
	if data.ScreensaverEndTime.String != "" && data.ScreensaverStartTime.String != "" {
		eq = eq.And(builder.Eq{
			"screensaver_start_time": data.ScreensaverStartTime,
			"screensaver_end_time":   data.ScreensaverEndTime,
			"screensaver_wait_time":  data.ScreenSaverWaitTime,
		})
	} else if data.ToppingEndTime.String != "" && data.ToppingStartTime.String != "" {
		eq = eq.And(builder.Eq{
			"topping_start_time": data.ToppingStartTime,
			"topping_end_time":   data.ToppingEndTime,
		})
	} else if data.PublishEndTime.String != "" && data.PublishStartTime.String != "" {
		eq = eq.And(builder.Eq{
			"publish_start_time": data.PublishStartTime,
			"publish_end_time":   data.PublishEndTime,
		})
	} else if data.LockScreenStartTime.String != "" && data.LockScreenEndTime.String != "" {
		eq = eq.And(builder.Eq{
			"lock_screen_start_time": data.LockScreenStartTime,
			"lock_screen_end_time":   data.LockScreenEndTime,
		})
	}
	we := builder.And(builder.Eq{"deleted": false})
	if data.DeviceId > 0 {
		we = we.And(builder.Eq{"device_id": data.DeviceId})
	}
	if data.PhotoId > 0 {
		we = we.And(builder.Eq{"photo_id": data.PhotoId})
	}
	toSQL, i, err := builder.Postgres().Update(eq).From(m.Table).Where(we).ToSQL()
	if err != nil {
		return err
	}
	_, err = m.ExecNoCache(toSQL, i...)
	return err
}

// ClearTimeByPhotoId 根据资源Id清除时间
func (m *DevicePhotoModel) ClearTimeByPhotoId(mod int, photoId int64) error {
	upTime := time.Now().Format("2006-01-02 15:04:05")
	eq := builder.And(builder.Eq{"updated_time": upTime})
	switch mod {
	case common.Topping:
		eq = eq.And(builder.Eq{
			"topping_start_time": nil,
			"topping_end_time":   nil,
		})
	case common.Publish:
		eq = eq.And(builder.Eq{
			"publish_start_time": nil,
			"publish_end_time":   nil,
		})
	case common.LockScreen:
		eq = eq.And(builder.Eq{
			"lock_screen_start_time": nil,
			"lock_screen_end_time":   nil,
		})
	case common.ScreenSaver:
		eq = eq.And(builder.Eq{
			"screensaver_start_time": nil,
			"screensaver_end_time":   sql.NullString{},
			"screensaver_wait_time":  0,
		})
	}
	toSQL, i, err := builder.Postgres().Update(eq).From(m.Table).Where(builder.Eq{
		"photo_id": photoId,
	}).ToSQL()
	fmt.Println(toSQL, eq)
	if err != nil {
		logx.Errorf("device clearTimeByPhotoId update err:%s", err.Error())
		return err
	}
	_, err = m.ExecNoCache(toSQL, i...)
	if err != nil {
		logx.Errorf("device clearTimeByPhotoId update err:%s", err.Error())
	}
	return err
}
func (m *DevicePhotoModel) Delete(id int64) error {

	photowallPhotosIdKey := fmt.Sprintf("%s%v", cachePhotowallPhotosIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where id = $1", m.Table)
		return conn.Exec(query, id)
	}, photowallPhotosIdKey)
	return err
}

func (m *DevicePhotoModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cachePhotowallPhotosIdPrefix, primary)
}

func (m *DevicePhotoModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1", photosFieldNames, m.Table)
	return conn.QueryRow(v, query, primary)
}
