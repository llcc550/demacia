package file

import (
	"encoding/json"
	"path"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
)

type (
	UpdateReq struct {
		Cache     *redis.Redis
		LocalPath string
		FileName  string
		OrgId     int64
		IsTmp     bool
	}
	huaweiObsConfig struct {
		Endpoint        string // obs.cn-east-3.myhuaweicloud.com
		AccessKeyId     string // 0LYNHSVNLVSB3IBDTXFN
		AccessKeySecret string // HahQCGzvvNN0kfEFOeoSNKVHx9WoU9SlU7leWQyx
		BucketName      string // u-test
		ObjectUrl       string // http://u-test.obs.cn-east-3.myhuaweicloud.com
	}
)

const huaweiConfigKey = "cache:upload:huawei"

func Upload(req *UpdateReq) (string, error) {
	data, err := req.Cache.Get(huaweiConfigKey)
	if err != nil {
		return "", err
	}
	var config huaweiObsConfig
	err = json.Unmarshal([]byte(data), &config)
	if err != nil {
		return "", err
	}
	uploadDir := "sourceFile"
	if req.IsTmp {
		uploadDir = "tmp"
	} else if req.OrgId != 0 {
		uploadDir = strconv.FormatInt(req.OrgId, 10)
	}
	uploadDir += time.Now().Format("/2006/01/")
	if req.FileName == "" {
		req.FileName = uuid.New().String() + path.Ext(req.LocalPath)
	}
	remoteUrl := uploadDir + req.FileName
	obsClient, err := obs.New(config.AccessKeyId, config.AccessKeySecret, config.Endpoint)
	if err != nil {
		return "", err
	}
	input := &obs.PutFileInput{}
	input.Bucket = config.BucketName
	input.Key = remoteUrl
	input.SourceFile = req.LocalPath
	_, err = obsClient.PutFile(input)
	if err != nil {
		return "", err
	}
	return config.ObjectUrl + "/" + remoteUrl, nil
}
