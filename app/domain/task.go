package domain

const TableNameTask = "task"

type Task struct {
}

func (Task) TableName() string { return TableNameTask }
