package logic

import (
	"context"
	"database/sql"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/position/api/internal/svc"
	"demacia/datacenter/position/api/internal/types"
	"demacia/datacenter/position/errors"
	"demacia/datacenter/position/model"
	"demacia/service/websocket/rpc/websocket"
	"fmt"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

type GeneratePositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGeneratePositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) GeneratePositionLogic {
	return GeneratePositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GeneratePositionLogic) GeneratePosition(req types.GeneratePositionReq) (*types.SuccessReply, error) {

	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return &types.SuccessReply{Success: false}, errlist.NoAuth
	}

	threading.GoSafe(func() {
		l.generatePosition(req.WebsocketUuid, oid)
	})
	return &types.SuccessReply{Success: true}, nil
}

func (l *GeneratePositionLogic) generatePosition(uuid string, oid int64) {
	fn := func(err *httpx.CodeError, msg ...interface{}) {
		_, _ = l.svcCtx.WebsocketRpc.Push(context.Background(), &websocket.Request{
			Key:  uuid,
			Code: int64(err.Code()),
			Msg:  fmt.Sprintf(err.Error(), msg...),
		})
	}

	classList, err := l.svcCtx.ClassRpc.ListByOrgId(context.Background(), &class.OrgIdReq{OrgId: oid})
	if err != nil {
		fn(errlist.Unknown)
	}

	for _, classInfo := range classList.List {
		_, err := l.svcCtx.PositionModel.SelectByClassId(classInfo.Id)
		if err != nil && err == sql.ErrNoRows {
			positionVal, err := l.svcCtx.PositionModel.SelectByClassFullName(classInfo.FullName, oid)
			if err != nil && err == sql.ErrNoRows {
				_, err = l.svcCtx.PositionModel.InsertPosition(&model.Position{
					Oid:          oid,
					PositionName: classInfo.FullName,
					ClassId:      classInfo.Id,
					ClassName:    classInfo.FullName,
				})
				if err != nil {
					fn(errors.PositionGenerateErr, classInfo.FullName)
				}
				fn(errors.PositionGenerateSuc, classInfo.FullName)
			} else if positionVal.PositionName == classInfo.FullName {
				if positionVal.ClassId != 0 {
					continue
				} else {
					err := l.svcCtx.PositionModel.UpdatePosition(&model.Position{
						Id:           positionVal.Id,
						Oid:          oid,
						PositionName: classInfo.FullName,
						ClassId:      classInfo.Id,
						ClassName:    classInfo.FullName,
					})
					if err != nil {
						fn(errors.PositionGenerateErr, classInfo.FullName)
					}
					fn(errors.PositionGenerateSuc, classInfo.FullName)
				}
			}
		}
	}

}
