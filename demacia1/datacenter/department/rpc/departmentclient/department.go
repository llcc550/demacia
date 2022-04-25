// Code generated by goctl. DO NOT EDIT!
// Source: department.proto

package departmentclient

import (
	"context"

	"demacia/datacenter/department/rpc/department"

	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	DepartmentIdReq            = department.DepartmentIdReq
	DepartmentIdsResp          = department.DepartmentIdsResp
	DepartmentInfo             = department.DepartmentInfo
	DepartmentListResp         = department.DepartmentListResp
	DepartmentMember           = department.DepartmentMember
	DepartmentMembersResp      = department.DepartmentMembersResp
	MemberIdsResp              = department.MemberIdsResp
	OrgIdAndDepartmentIdReq    = department.OrgIdAndDepartmentIdReq
	OrgIdAndDepartmentTitleReq = department.OrgIdAndDepartmentTitleReq
	OrgIdAndMemberIdReq        = department.OrgIdAndMemberIdReq
	OrgIdReq                   = department.OrgIdReq

	Department interface {
		//  根据id获取部门详情
		GetDepartmentById(ctx context.Context, in *DepartmentIdReq, opts ...grpc.CallOption) (*DepartmentInfo, error)
		//  根据机构ID和部门名称获取部门信息
		GetDepartmentByOrgIdAndDepartmentTitle(ctx context.Context, in *OrgIdAndDepartmentTitleReq, opts ...grpc.CallOption) (*DepartmentInfo, error)
		//  根据机构ID和部门ID获取部门下的人员ID列表
		GetMemberIdsByOrgIdAndDepartmentId(ctx context.Context, in *OrgIdAndDepartmentIdReq, opts ...grpc.CallOption) (*MemberIdsResp, error)
		//  根据机构ID和人员ID获取该人员所在的部门ID列表
		GetDepartmentIdsByOrgIdAndMemberId(ctx context.Context, in *OrgIdAndMemberIdReq, opts ...grpc.CallOption) (*DepartmentIdsResp, error)
		//  根据机构ID获取所有的部门
		GetDepartmentsByOrgId(ctx context.Context, in *OrgIdReq, opts ...grpc.CallOption) (*DepartmentListResp, error)
		//  根据机构ID获取所有的部门人员关系
		GetDepartmentMemberRelationByOrgId(ctx context.Context, in *OrgIdReq, opts ...grpc.CallOption) (*DepartmentMembersResp, error)
	}

	defaultDepartment struct {
		cli zrpc.Client
	}
)

func NewDepartment(cli zrpc.Client) Department {
	return &defaultDepartment{
		cli: cli,
	}
}

//  根据id获取部门详情
func (m *defaultDepartment) GetDepartmentById(ctx context.Context, in *DepartmentIdReq, opts ...grpc.CallOption) (*DepartmentInfo, error) {
	client := department.NewDepartmentClient(m.cli.Conn())
	return client.GetDepartmentById(ctx, in, opts...)
}

//  根据机构ID和部门名称获取部门信息
func (m *defaultDepartment) GetDepartmentByOrgIdAndDepartmentTitle(ctx context.Context, in *OrgIdAndDepartmentTitleReq, opts ...grpc.CallOption) (*DepartmentInfo, error) {
	client := department.NewDepartmentClient(m.cli.Conn())
	return client.GetDepartmentByOrgIdAndDepartmentTitle(ctx, in, opts...)
}

//  根据机构ID和部门ID获取部门下的人员ID列表
func (m *defaultDepartment) GetMemberIdsByOrgIdAndDepartmentId(ctx context.Context, in *OrgIdAndDepartmentIdReq, opts ...grpc.CallOption) (*MemberIdsResp, error) {
	client := department.NewDepartmentClient(m.cli.Conn())
	return client.GetMemberIdsByOrgIdAndDepartmentId(ctx, in, opts...)
}

//  根据机构ID和人员ID获取该人员所在的部门ID列表
func (m *defaultDepartment) GetDepartmentIdsByOrgIdAndMemberId(ctx context.Context, in *OrgIdAndMemberIdReq, opts ...grpc.CallOption) (*DepartmentIdsResp, error) {
	client := department.NewDepartmentClient(m.cli.Conn())
	return client.GetDepartmentIdsByOrgIdAndMemberId(ctx, in, opts...)
}

//  根据机构ID获取所有的部门
func (m *defaultDepartment) GetDepartmentsByOrgId(ctx context.Context, in *OrgIdReq, opts ...grpc.CallOption) (*DepartmentListResp, error) {
	client := department.NewDepartmentClient(m.cli.Conn())
	return client.GetDepartmentsByOrgId(ctx, in, opts...)
}

//  根据机构ID获取所有的部门人员关系
func (m *defaultDepartment) GetDepartmentMemberRelationByOrgId(ctx context.Context, in *OrgIdReq, opts ...grpc.CallOption) (*DepartmentMembersResp, error) {
	client := department.NewDepartmentClient(m.cli.Conn())
	return client.GetDepartmentMemberRelationByOrgId(ctx, in, opts...)
}
