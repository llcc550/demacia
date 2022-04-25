package logic

import (
	"context"
	"fmt"
	"sort"

	"demacia/common/baseauth"
	"demacia/datacenter/department/rpc/department"
	"demacia/datacenter/member/api/internal/svc"
	"demacia/datacenter/member/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type ChoiceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChoiceLogic(ctx context.Context, svcCtx *svc.ServiceContext) ChoiceLogic {
	return ChoiceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChoiceLogic) Choice() (resp *types.MemberChoiceRes, err error) {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, err
	}
	// 获取所有部门
	_, depMap := l.GetAllDepartment(orgId)
	//fmt.Println("depMap:", depMap)

	// 追加未分组的部门
	depMap[0] = &types.DepartmentInfo{
		DepartmentId:    0,
		DepartmentTitle: "未分组",
		Children:        nil,
		Count:           0,
	}
	// 获取所有用户
	memberIds, memberMap := l.GetAllMember(orgId)

	// 获取成员部门关联关系
	md, dm := l.GetDepartmentMember(orgId)
	fmt.Println("md:", md)
	// 获取未参与任何部门的成员 ids
	n := l.GetNotIntAnyDepartmentMember(memberIds, md)
	//归类
	list, _ := l.GetMemberList(dm, memberMap)

	// 未分组
	for _, ntv := range n {
		depMap[0].Children = append(depMap[0].Children, memberMap[ntv])
		dm[0] = append(dm[0], ntv)
		list[0] = append(list[0], memberMap[ntv])
	}

	var res = types.MemberChoiceRes{}
	// 循环部门map
	for dk, dv := range depMap {
		di := types.DepartmentInfo{}
		di.DepartmentTitle = dv.DepartmentTitle
		di.DepartmentId = dv.DepartmentId

		// 赋值
		if _, ok := dm[dk]; ok {

			di.Children = list[dk]
			di.Count = int64(len(list[dk]))
			// 修改pid

		} else {
			di.Children = []types.Member{}
			di.Count = 0
		}

		res.DepartmentList = append(res.DepartmentList, di)
	}
	// 结构体切片排序
	sort.Slice(res.DepartmentList, func(i, j int) bool {
		return res.DepartmentList[i].DepartmentId < res.DepartmentList[j].DepartmentId
	})

	return &res, nil
}

// GetAllDepartment 获取所有部门
func (l *ChoiceLogic) GetAllDepartment(orgId int64) ([]int64, map[int64]*types.DepartmentInfo) {
	// 获取所有部门
	departmentList, _ := l.svcCtx.DepartmentRpc.GetDepartmentsByOrgId(l.ctx, &department.OrgIdReq{OrgId: orgId})

	depMap := make(map[int64]*types.DepartmentInfo)
	var depIds []int64
	// 循环部门列表
	for _, v := range departmentList.Departments {
		depMap[v.DepartmentId] = &types.DepartmentInfo{
			DepartmentId:    v.DepartmentId,
			DepartmentTitle: v.DepartmentTitle,
			Children:        []types.Member{},
			Count:           0,
		}
		depIds = append(depIds, v.DepartmentId)
	}
	return depIds, depMap

}

// GetAllMember 获取机构下所有成员
func (l *ChoiceLogic) GetAllMember(orgId int64) ([]int64, map[int64]types.Member) {
	// 查询数据库
	members, _ := l.svcCtx.MemberModel.FindListByOrgId(orgId)
	memberMap := make(map[int64]types.Member)
	var memberIds []int64
	for _, m := range members {
		memberMap[m.Id] = types.Member{
			MemberId:   m.Id,
			UserName:   m.UserName,
			TrueName:   m.TrueName,
			Mobile:     m.Mobile,
			Status:     m.Status,
			Face:       m.Face,
			FaceStatus: m.FaceStatus,
		}
		memberIds = append(memberIds, m.Id)
	}
	return memberIds, memberMap
}

// GetDepartmentMember 获取成员部门关联关系
func (l *ChoiceLogic) GetDepartmentMember(orgId int64) (md, dm map[int64][]int64) {
	dml, _ := l.svcCtx.DepartmentRpc.GetDepartmentMemberRelationByOrgId(l.ctx, &department.OrgIdReq{OrgId: orgId})

	// 人员对应 多个部门
	md = map[int64][]int64{}
	// 部门对应 多个人员
	dm = map[int64][]int64{}
	for _, v := range dml.DepartmentMembers {
		md[v.MemberId] = append(md[v.MemberId], v.DepartmentId)
		dm[v.DepartmentId] = append(dm[v.DepartmentId], v.MemberId)
	}
	return md, dm
}

// GetNotIntAnyDepartmentMember 获取所有未在任何分组下的成员
func (l *ChoiceLogic) GetNotIntAnyDepartmentMember(memberIds []int64, md map[int64][]int64) []int64 {
	var nt []int64
	for _, v := range memberIds {
		if _, ok := md[v]; !ok {
			nt = append(nt, v)
		}
	}
	return nt
}

// GetMemberList 归类
func (l *ChoiceLogic) GetMemberList(dm map[int64][]int64, members map[int64]types.Member) (map[int64][]types.Member, error) {
	dml := map[int64][]types.Member{}
	for dmk, dmv := range dm {
		for _, dmvV := range dmv {
			dml[dmk] = append(dml[dmk], types.Member{
				MemberId:   members[dmvV].MemberId,
				UserName:   members[dmvV].UserName,
				TrueName:   members[dmvV].TrueName,
				Mobile:     members[dmvV].Mobile,
				Status:     members[dmvV].Status,
				Face:       members[dmvV].Face,
				FaceStatus: members[dmvV].FaceStatus,
				Pid:        dmk,
			})
		}
	}
	return dml, nil
}
