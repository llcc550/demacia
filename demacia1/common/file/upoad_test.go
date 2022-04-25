package file

import (
	"fmt"
	"testing"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
)

func TestUpload(t *testing.T) {
	cache := redis.New("122.112.230.215:6379")
	req := UpdateReq{
		Cache:     cache,
		LocalPath: "test.txt",
		FileName:  "test.txt",
		OrgId:     0,
		IsTmp:     false,
	}
	url, err := Upload(&req)
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	fmt.Println(url)
}
