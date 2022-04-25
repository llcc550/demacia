package model

import (
	"context"

	"github.com/globalsign/mgo/bson"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/mongo"
)

type LogModel interface {
	Insert(ctx context.Context, data *Log) error
	FindOne(ctx context.Context, id string) (*Log, error)
	Update(ctx context.Context, data *Log) error
	Delete(ctx context.Context, id string) error
}

type defaultLogModel struct {
	*mongo.Model
}

func NewLogModel(url, collection string) LogModel {
	return &defaultLogModel{
		Model: mongo.MustNewModel(url, collection),
	}
}

func (m *defaultLogModel) Insert(ctx context.Context, data *Log) error {
	if !data.ID.Valid() {
		data.ID = bson.NewObjectId()
	}

	session, err := m.TakeSession()
	if err != nil {
		return err
	}

	defer m.PutSession(session)
	return m.GetCollection(session).Insert(data)
}

func (m *defaultLogModel) FindOne(ctx context.Context, id string) (*Log, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, ErrInvalidObjectId
	}

	session, err := m.TakeSession()
	if err != nil {
		return nil, err
	}

	defer m.PutSession(session)
	var data Log

	err = m.GetCollection(session).FindId(bson.ObjectIdHex(id)).One(&data)
	switch err {
	case nil:
		return &data, nil
	case mongo.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultLogModel) Update(ctx context.Context, data *Log) error {
	session, err := m.TakeSession()
	if err != nil {
		return err
	}

	defer m.PutSession(session)

	return m.GetCollection(session).UpdateId(data.ID, data)
}

func (m *defaultLogModel) Delete(ctx context.Context, id string) error {
	session, err := m.TakeSession()
	if err != nil {
		return err
	}

	defer m.PutSession(session)

	return m.GetCollection(session).RemoveId(bson.ObjectIdHex(id))
}
