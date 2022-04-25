package datacenter

import "encoding/json"

const Kafka = "datacenter"

const (
	Organization = "organization"
	Stage        = "stage"
	Grade        = "grade"
	Class        = "class"
	Device       = "device"
	Member       = "member"
	Department   = "department"
	Student      = "student"
	Parent       = "parent"
	CourseTable  = "coursetable"
	Position     = "position"
	TimeSwitch   = "timeswitch"
	Subject      = "subject"
)

const (
	Create = int8(1)
	Add    = int8(1) // todo: 废除，使用Create常量
	Update = int8(2)
	Delete = int8(3)
)

type (
	Message struct {
		Topic    string `json:"topic"`     // 数据主题。
		ObjectId int64  `json:"object_id"` // 资源数据ID
		Action   int8   `json:"action"`    // 1:新增,2:修改,3:删除
		Data     string `json:"data"`      // todo: 废除
	}
)

func Unmarshal(data string) *Message {
	var message Message
	_ = json.Unmarshal([]byte(data), &message)
	return &message
}

// 以下废除

type (
	OrganizationData struct {
		Title string `json:"title"`
	}
	DepartmentData struct {
		Title string `json:"title"`
	}
	StudentData struct {
		StudentName    string `json:"student_name"`
		UserName       string `json:"user_name"`
		OrganizationId int64  `json:"organization_id"`
		ClassId        int64  `json:"class_id"`
		Sex            int8   `json:"sex"`
	}
	ParentData struct {
		ParentName  string        `json:"parent_name"`
		UserName    string        `json:"user_name"`
		Mobile      string        `json:"mobile"`
		StudentInfo []StudentInfo `json:"student_info"`
	}
	StudentInfo struct {
		StudentId int64 `json:"student_id"`
		Relation  int8  `json:"relation"`
	}
)

func Marshal(topic string, objectId int64, action int8, data interface{}) string {
	ds := ""
	if data != nil {
		d, _ := json.Marshal(data)
		ds = string(d)
	}
	s, _ := json.Marshal(Message{
		Topic:    topic,
		ObjectId: objectId,
		Data:     ds,
		Action:   action,
	})
	return string(s)
}
