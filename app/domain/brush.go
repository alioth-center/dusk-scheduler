package domain

const TableNameBrush = "brush"

type Brush struct{}

func (Brush) TableName() string { return TableNameBrush }
