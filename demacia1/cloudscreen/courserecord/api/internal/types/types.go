// Code generated by goctl. DO NOT EDIT.
package types

type CourseRecordAddReq struct {
	UserId   int64  `json:"user_id"`
	UserType int8   `json:"user_type"`
	SignTime string `json:"sign_time"`
	Photo    string `json:"photo"`
}

type SuccessReply struct {
	Success bool `json:"success"`
}

type RecordsReq struct {
	StartDate string `json:"start_date,optional"`
	EndDate   string `json:"end_date,optional"`
	Status    int8   `json:"status,default=-1"`
	Truename  string `json:"truename,optional"`
	UserType  int8   `json:"user_type,default=-1"`
	Page      int    `json:"page,default=1"`
	Limit     int    `json:"limit,default=10"`
}

type RecordsReply struct {
	Count      int           `json:"count"`
	RecordList []*RecordInfo `json:"record_list"`
}

type RecordInfo struct {
	Name         string `json:"name"`
	ClassName    string `json:"class_name,omitempty"`
	Date         string `json:"date"`
	SignTime     string `json:"sign_time"`
	Note         string `json:"note"`
	SubjectName  string `json:"subject_name"`
	PositionName string `json:"position_name,omitempty"`
	Status       int8   `json:"status"`
	Photo        string `json:"photo"`
}

type CourseRecordConfigReply struct {
	Enable         int            `json:"enable"`
	SignPerson     int            `json:"sign_person"`
	SignBeforeTime int            `json:"sign_before_time"`
	SignHoliday    int            `json:"sign_holiday"`
	SpecialDates   []*SpecialDate `json:"special_dates"`
}

type SpecialDate struct {
	Year int    `json:"year"`
	Date string `json:"date"`
	Type int    `json:"type"`
}

type CourseRecordConfigSaveReq struct {
	Enable         int            `json:"enable"`
	SignPerson     int            `json:"sign_person"`
	SignBeforeTime int            `json:"sign_before_time"`
	SignHoliday    int            `json:"sign_holiday"`
	SpecialDates   []*SpecialDate `json:"special_dates"`
}

type CourseRecordCountReq struct {
	Page        int    `json:"page,default=1"`
	Limit       int    `json:"limit,default=10"`
	Type        int    `json:"type"`
	StageId     int64  `json:"stage_id,optional"`
	GradeId     int64  `json:"grade_id,optional"`
	ClassId     int64  `json:"class_id,optional"`
	StartDate   string `json:"start_date,optional"`
	EndDate     string `json:"end_date,optional"`
	StudentName string `json:"student_name,optional"`
}

type StudentRecord struct {
	ClassName   string `json:"class_name"`
	StudentName string `json:"student_name"`
	Attendance  string `json:"attendance"`
	ShouldCount int    `json:"should_count"`
	NormalCount int    `json:"normal_count"`
	LateCount   int    `json:"late_count"`
	LackCount   int    `json:"lack_count"`
}

type StudentRecordReply struct {
	Count          int              `json:"count"`
	AllAttendance  string           `json:"all_attendance"`
	AllShouldCount int              `json:"all_should_count"`
	AllNormalCount int              `json:"all_normal_count"`
	AllLateCount   int              `json:"all_late_count"`
	AllLackCount   int              `json:"all_lack_count"`
	StudentRecords []*StudentRecord `json:"student_records"`
}

type InitReq struct {
	ClassId int64 `json:"class_id"`
}

type TeacherInfo struct {
	TeacherName string `json:"teacher_name"`
	Status      int    `json:"status"`
	Photo       string `json:"photo"`
}

type StudentInfo struct {
	StudentName string `json:"student_name"`
	Status      int    `json:"status"`
	Photo       string `json:"photo"`
}

type CourseInfo struct {
	SubjectName string `json:"subject_name"`
	Note        string `json:"note"`
}

type ClassCourseRecordInfoReply struct {
	StartTime   string         `json:"start_time"`
	EndTime     string         `json:"end_time"`
	NormalCount int            `json:"normal_count"`
	LateCount   int            `json:"late_count"`
	LackCount   int            `json:"lack_count"`
	Course      *CourseInfo    `json:"course"`
	Teacher     *TeacherInfo   `json:"teacher,omitempty"`
	Students    []*StudentInfo `json:"students"`
}

type ClassCourseRecordReq struct {
	Page        int    `json:"page,default=1"`
	Limit       int    `json:"limit,default=10"`
	ClassId     int64  `json:"class_id"`
	Truename    string `json:"truename,optional"`
	SubjectName string `json:"subject_name,optional"`
	QueryDate   string `json:"query_date,optional"`
	Status      int8   `json:"status,default=-1"`
}

type ClassCourseRecordReply struct {
	Count              int           `json:"count"`
	ClassCourseRecords []*RecordInfo `json:"class_course_reocods"`
}

type CourseRecordReq struct {
	Page      int   `json:"page,default=1"`
	Limit     int   `json:"limit,default=10"`
	Type      int   `json:"type"`
	StudentId int64 `json:"student_id,optional"`
	MemberId  int64 `json:"member_id,optional"`
}

type CourseRecordReply struct {
	Count         int           `json:"count"`
	CourseRecords []*RecordInfo `json:"course_reocods"`
}

type PhotoReq struct {
	Photo string `json:"photo"`
}
