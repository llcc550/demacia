package main

import (
	"context"
	"flag"
	"time"

	"demacia/datacenter/department/rpc/departmentclient"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

func main() {
	flag.Parse()

	client, err := zrpc.NewClient(zrpc.NewDirectClientConf([]string{"127.0.0.1:2100"}, "app", "token"))
	if err != nil {
		logx.Error(err)
	}
	cli := departmentclient.NewDepartment(client)

	res1, err1 := cli.GetDepartmentById(context.Background(), &departmentclient.DepartmentIdReq{
		DepartmentId: 1,
	})
	if err1 != nil {
		logx.Error(err1.Error())
	} else {
		logx.Info(res1)
	}
	time.Sleep(time.Second)

	res2, err2 := cli.GetDepartmentByOrgIdAndDepartmentTitle(context.Background(), &departmentclient.OrgIdAndDepartmentTitleReq{
		OrgId:           2,
		DepartmentTitle: "12",
	})
	if err2 != nil {
		logx.Error(err2.Error())
	} else {
		logx.Info(res2)
	}
	time.Sleep(time.Second)

	res3, err3 := cli.GetMemberIdsByOrgIdAndDepartmentId(context.Background(), &departmentclient.OrgIdAndDepartmentIdReq{
		OrgId:        2,
		DepartmentId: 1,
	})
	if err3 != nil {
		logx.Error(err3.Error())
	} else {
		logx.Info(res3)
	}
	time.Sleep(time.Second)

	res4, err4 := cli.GetDepartmentsByOrgId(context.Background(), &departmentclient.OrgIdReq{
		OrgId: 2,
	})
	if err4 != nil {
		logx.Error(err4.Error())
	} else {
		logx.Info(res4)
	}
	time.Sleep(time.Second)

	res5, err5 := cli.GetDepartmentMemberRelationByOrgId(context.Background(), &departmentclient.OrgIdReq{
		OrgId: 2,
	})
	if err4 != nil {
		logx.Error(err5.Error())
	} else {
		logx.Info(res5)
	}
	time.Sleep(time.Second)
}
