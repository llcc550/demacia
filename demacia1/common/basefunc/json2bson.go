package basefunc

import (
	"encoding/json"

	"github.com/globalsign/mgo/bson"
)

func Json2Bson(j string) string {
	j = RemoveSpace(j)
	var ji interface{}
	_ = json.Unmarshal([]byte(j), &ji)
	b, _ := bson.Marshal(ji)
	return string(b)
}
