# Demacia: Now and forever.

## tips

1. postgres生产model文件

```shell
goctl model pg datasource -url="postgres://postgres:dinghao8888@122.112.230.215:5432/demacia?sslmode=disable" -schema="card"  -table="card" -dir ./model -cache=true
```

2. 测试环境基础服务

- [x] 数据库
    - postgres:dinghao8888@122.112.230.215:5432/demacia
- [x] websocket
    - 124.71.163.192:32301
    - 124.71.163.192:32302
    - 124.71.163.192:32303
- [x] kafka
    - 124.71.163.192:32420
    - 124.71.163.192:32421
    - 124.71.163.192:32422
- [x] mongodb
    - 124.71.163.192:32500
- [x] redis
    - 124.71.163.192:32600
- [x] beanstalkd
    - 124.71.163.192:32700
    - 124.71.163.192:32701
    - 124.71.163.192:32702
- [x] dtm:
    - 124.71.163.192:31790

3. 测试环境RPC端口列表

> ip 124.71.163.192

| rpc             | port  | 备注        |
|-----------------|-------|-----------|
| databus         | 32000 | 数据总线      |
| organization    | 32001 | 机构        |
| common          | 32002 | 通用        |
| member          | 32003 | 教师        |
| class           | 32004 | 班级        |
| websocket       | 32005 | websocket |
| user            | 32006 | 中控用户      |
| card            | 32007 | 卡号        |
| device          | 32008 | 设备        |
| student         | 32009 | 学生        |
| subject         | 32010 | 学科        |
| position        | 32011 | 地点        |
| department      | 32012 | 部门        |
| coursetable-rpc | 32013 | 课程表       |
| tags-rpc        | 32014 | 云盘标签      |
| capacity-rpc    | 32015 | 云盘权限      |

4.obs使用
> 以下key:value保存到本地开发环境redis。或本地直接使用测试环境redis：[124.71.163.192:32600]

* cache:upload:huawei

```json
{
  "Endpoint": "obs.cn-east-3.myhuaweicloud.com",
  "AccessKeyId": "0LYNHSVNLVSB3IBDTXFN",
  "AccessKeySecret": "HahQCGzvvNN0kfEFOeoSNKVHx9WoU9SlU7leWQyx",
  "BucketName": "u-test",
  "ObjectUrl": "http://u-test.obs.cn-east-3.myhuaweicloud.com"
}
```