package types

type (
	ChildrenList    []*Children
	DeviceGroupList []*DeviceGroup
)

func (t ChildrenList) Len() int           { return len(t) }
func (t ChildrenList) Swap(i, j int)      { (t)[i], (t)[j] = (t)[j], (t)[i] }
func (t ChildrenList) Less(i, j int) bool { return t[i].Id < t[j].Id }

func (t DeviceGroupList) Len() int           { return len(t) }
func (t DeviceGroupList) Swap(i, j int)      { (t)[i], (t)[j] = (t)[j], (t)[i] }
func (t DeviceGroupList) Less(i, j int) bool { return t[i].Id < t[j].Id }
