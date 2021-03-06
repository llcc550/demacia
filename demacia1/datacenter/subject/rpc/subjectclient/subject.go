// Code generated by goctl. DO NOT EDIT!
// Source: subject.proto

package subjectclient

import (
	"context"

	"demacia/datacenter/subject/rpc/subject"

	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	IdReq       = subject.IdReq
	NullResp    = subject.NullResp
	SubjectInfo = subject.SubjectInfo

	Subject interface {
		GetSubjectById(ctx context.Context, in *IdReq, opts ...grpc.CallOption) (*SubjectInfo, error)
	}

	defaultSubject struct {
		cli zrpc.Client
	}
)

func NewSubject(cli zrpc.Client) Subject {
	return &defaultSubject{
		cli: cli,
	}
}

func (m *defaultSubject) GetSubjectById(ctx context.Context, in *IdReq, opts ...grpc.CallOption) (*SubjectInfo, error) {
	client := subject.NewSubjectClient(m.cli.Conn())
	return client.GetSubjectById(ctx, in, opts...)
}
