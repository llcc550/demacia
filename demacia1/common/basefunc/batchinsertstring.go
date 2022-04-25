package basefunc

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func BatchInsertString(tableName string, params interface{}) (string, []interface{}, error) {
	s := reflect.ValueOf(params)
	if s.Kind() != reflect.Slice {
		return "", nil, errors.New("参数错误")
	}
	typeOfStruct := reflect.TypeOf(s.Index(0).Interface()).Elem()
	rowsExpectIdAutoSet := ""
	var inserts []string
	n := 1
	var args []interface{}
	for i := 0; i < s.Len(); i++ {
		valueOfData := reflect.ValueOf(s.Index(i).Interface()).Elem()
		num := valueOfData.NumField()
		insert := ""
		for j := 0; j < num; j++ {
			if typeOfStruct.Field(j).Tag.Get("db") == "id" {
				continue
			}
			if i == 0 {
				rowsExpectIdAutoSet += typeOfStruct.Field(j).Tag.Get("db") + ","
			}
			insert = fmt.Sprintf("%s,$%d", insert, n)
			n++
			f := valueOfData.Field(j)
			switch f.Type().String() {
			case "uint":
				fallthrough
			case "uint8":
				fallthrough
			case "uint16":
				fallthrough
			case "uint32":
				fallthrough
			case "uint64":
				args = append(args, f.Uint())
			case "int":
				fallthrough
			case "int8":
				fallthrough
			case "int16":
				fallthrough
			case "int32":
				fallthrough
			case "int64":
				args = append(args, f.Int())
			case "float32":
				fallthrough
			case "float64":
				args = append(args, f.Float())
			case "time.Time":
				value := f.Interface().(time.Time)
				if value.Unix() < 0 {
					args = append(args, time.Now())
				} else {
					args = append(args, value)
				}
			case "sql.NullTime":
				value := f.Interface().(sql.NullTime)
				if value.Valid {
					args = append(args, value.Time)
				} else {
					args = append(args, nil)
				}
			case "sql.NullString":
				value := f.Interface().(sql.NullString)
				if value.Valid {
					args = append(args, value.String)
				} else {
					args = append(args, nil)
				}
			case "string":
				args = append(args, f.String())
			default:
				return "", nil, errors.New("数据中有尚未支持的数据结构" + f.Type().String())
			}
		}
		inserts = append(inserts, "("+strings.Trim(insert, ",")+")")
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", tableName, strings.Trim(rowsExpectIdAutoSet, ","), strings.Join(inserts, ","))
	return query, args, nil
}
