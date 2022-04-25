package model

//go:generate goctl model mongo -t Log
import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type Log struct {
	ID         bson.ObjectId `bson:"_id"`
	Ip         string        `bson:"ip"`
	Route      string        `bson:"route"`
	Jwt        interface{}   `bson:"jwt"`
	Req        interface{}   `bson:"req"`
	CreateTime time.Time     `bson:"create_time"`
}
