package domain

const TaskNameOutcome = "outcome"

type Outcome struct {
}

func (Outcome) TableName() string { return TaskNameOutcome }
