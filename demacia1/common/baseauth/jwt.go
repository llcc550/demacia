package baseauth

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"demacia/common/errlist"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type AuthConfig struct {
	AccessSecret      string
	RefreshSecret     string `json:",optional"`
	PrevRefreshSecret string `json:",optional"`
	AccessExpire      int64  `json:",default=1209600"`
	RefreshExpire     int64  `json:",default=2419200"`
	RefreshAfter      int64  `json:",default=604800"`
}

const (
	OrgIdField    = "org_id"    // 机构ID int64
	MemberIdField = "member_id" // 人员ID int64
	RoleField     = "role"      // 身份 int64 1:教师，2：家长
	UserIdField   = "user_id"   // 用户ID int64
)

type (
	UserJwt struct {
		OrgId    int64 `json:"org_id"`
		MemberId int64 `json:"member_id"`
		Role     int64 `json:"role"`
		UserId   int64 `json:"user_id"`
	}
)

// GetOrgId 从jwt中获取org_id
func GetOrgId(r context.Context) (int64, error) {
	return getInt64FromRequest(r, OrgIdField)
}

// GetMemberId 从jwt中获取member_id
func GetMemberId(r context.Context) (int64, error) {
	return getInt64FromRequest(r, MemberIdField)
}

// GetUserId 从jwt中获取user_id
func GetUserId(r context.Context) (int64, error) {
	return getInt64FromRequest(r, UserIdField)
}

func GetUserJwt(r *http.Request) (*UserJwt, error) {
	userJwt := new(UserJwt)
	err := getJwt(r, userJwt)
	if err != nil {
		return nil, errlist.NoAuth
	}
	return userJwt, nil
}

func getJwt(r *http.Request, jwt interface{}) error {
	v := reflect.ValueOf(jwt).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		name := fieldInfo.Tag.Get("json")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		} else {
			name = strings.Split(name, ",")[0]
		}
		jwtValue := r.Context().Value(name)
		// json.Number类型转换为int64、Float64或string
		if value, ok := jwtValue.(json.Number); ok {
			if valueInt64, err := value.Int64(); err == nil {
				jwtValue = valueInt64
			} else if valueFloat64, err := value.Float64(); err == nil {
				jwtValue = valueFloat64
			} else {
				return errlist.NoAuth
			}
		} else if value, ok := jwtValue.(string); ok {
			jwtValue = value
		} else {
			continue
		}
		if reflect.ValueOf(jwtValue).Type() == v.FieldByName(fieldInfo.Name).Type() {
			v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(jwtValue))
		}
	}
	return nil
}

// 通用：从jwt中提取类型为string的数据
func getStringFromRequest(r context.Context, jwtKey string) (string, error) {
	field := r.Value(jwtKey)
	if field == nil {
		logx.Errorf("request not has %s", jwtKey)
		return "", errlist.NoAuth
	}
	jwtValue, ok := field.(string)
	if !ok {
		logx.Errorf("request with wrong type %s, param is: %+v", jwtKey, field)
	}
	if len(jwtValue) == 0 {
		logx.Errorf("request with %s: %s length is zero, param is: %+v, ", jwtKey, jwtValue, field)
		return "", errlist.NoAuth
	}
	return jwtValue, nil
}

// 通用：从jwt中提取类型为int64的数据
func getInt64FromRequest(r context.Context, jwtKey string) (int64, error) {
	field := r.Value(jwtKey)
	if field == nil {
		logx.Errorf("request not has %s", jwtKey)
		return 0, errlist.NoAuth
	}

	jwtJsonNumberValue, ok := field.(json.Number)
	if !ok {
		logx.Errorf("request with wrong type %s, param is: %+v", jwtKey, field)
		return 0, errlist.NoAuth
	}

	jwtValue, err := jwtJsonNumberValue.Int64()
	if err != nil {
		logx.Errorf("request with wrong type %s, param is: %+v", jwtKey, field)
	}

	if jwtValue == 0 {
		logx.Errorf("request with %s: %d  is 0, param is: %+v, ", jwtKey, jwtValue, field)
		return 0, errlist.NoAuth
	}

	return jwtValue, nil
}
