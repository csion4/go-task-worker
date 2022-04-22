package vo


type TaskVO struct {
	TaskCode string
	RecordId int
	Stages []map[int]map[string]string
}

