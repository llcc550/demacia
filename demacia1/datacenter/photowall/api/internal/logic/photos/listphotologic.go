package photos

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/photowall/api/internal/svc"
	"demacia/datacenter/photowall/api/internal/types"
	"demacia/datacenter/photowall/model"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"time"
)

type ListPhotoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPhotoLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListPhotoLogic {
	return ListPhotoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPhotoLogic) ListPhoto(req types.ListPhotoReq) (resp *types.ListPhoto, err error) {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, errlist.NoAuth
	}
	if req.FolderId == 0 && req.DeviceId == 0 || req.FolderId > 0 && req.DeviceId > 0 {
		return nil, errlist.InvalidParam
	}
	if req.FolderId != 0 {
		resp, err = l.ListPhotoByFolder(orgId, req)
	} else if req.DeviceId != 0 {
		resp, err = l.ListPhotoByDevice(orgId, req)
	}
	if err != nil {
		return nil, errlist.Unknown
	}
	return resp, nil
}

// ListPhotoByDevice 按设备获取列表
func (l *ListPhotoLogic) ListPhotoByDevice(orgId int64, req types.ListPhotoReq) (resp *types.ListPhoto, err error) {
	list, total, err := l.svcCtx.DevicePhotoModel.ListByDeviceId(&model.ListByIdReq{
		Id:    req.DeviceId,
		Title: req.Title,
		OrgId: orgId,
		Page:  req.Page,
		Limit: req.Limit,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.ListPhoto{
		List:  []types.PhotoInfo{},
		Total: 0,
	}
	resp.Total = total
	currentTimeUnix := time.Now().Unix()
	timeLayout := "2006-01-02T15:04:05Z"
	for _, v := range list {
		Info := types.PhotoInfo{
			Id:                v.PhotoId,
			Title:             v.DevicePhotoTitle,
			Url:               v.PhotoUrl,
			PublishStatus:     0,
			ScreenSaverStatus: 0,
			LockScreenStatus:  0,
			ToppingStatus:     0,
		}
		// 发布状态
		if v.PublishEndTime.String != "" {
			theTime, _ := time.ParseInLocation(timeLayout, v.PublishEndTime.String, time.Local)
			publishEndTimeUnix := theTime.Unix()
			if publishEndTimeUnix > currentTimeUnix {
				Info.PublishStatus = 1
			}
		}
		// 锁屏状态
		if v.LockScreenEndTime.String != "" {
			theTime, _ := time.ParseInLocation(timeLayout, v.LockScreenEndTime.String, time.Local)
			lockScreenEndTimeUnix := theTime.Unix()
			if lockScreenEndTimeUnix > currentTimeUnix {
				Info.LockScreenStatus = 1
			}
		}
		// 置顶状态
		if v.ToppingEndTime.String != "" {
			theTime, _ := time.ParseInLocation(timeLayout, v.ToppingEndTime.String, time.Local)
			toppingEndTimeUnix := theTime.Unix()
			if toppingEndTimeUnix > currentTimeUnix {
				Info.ToppingStatus = 1
			}
		}
		// 屏保状态
		if v.ScreensaverEndTime.String != "" {
			theTime, _ := time.ParseInLocation(timeLayout, v.ScreensaverEndTime.String, time.Local)
			screensaverEndTimeUnix := theTime.Unix()
			if screensaverEndTimeUnix > currentTimeUnix {
				Info.ScreenSaverStatus = 1
			}
		}
		resp.List = append(resp.List, Info)

	}

	// todo : 后续补上资源对于设备的，屏保，锁屏，置顶,发布状态
	return resp, nil
}

// ListPhotoByFolder 按相册获取资源列表
func (l *ListPhotoLogic) ListPhotoByFolder(orgId int64, req types.ListPhotoReq) (resp *types.ListPhoto, err error) {

	list, total, err := l.svcCtx.PhotosModel.List(&model.ListReq{
		PhotoFolderId: req.FolderId,
		Title:         req.Title,
		OrgId:         orgId,
		Page:          req.Page,
		Limit:         req.Limit,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.ListPhoto{
		List:  []types.PhotoInfo{},
		Total: 0,
	}
	resp.Total = total
	for _, v := range list {
		resp.List = append(resp.List, types.PhotoInfo{
			Id:                v.Id,
			Title:             v.Title,
			Url:               v.Url,
			CreatedTime:       v.CreatedTime,
			PublishStatus:     0,
			ScreenSaverStatus: 0,
			LockScreenStatus:  0,
			ToppingStatus:     0,
		})
	}
	// todo : 后续补上资源对于相册，屏保，锁屏，置顶,发布状态
	return resp, nil
}
