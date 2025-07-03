package domain

const TableNamePainter = "painter"

type Painter struct{}

func (Painter) TableName() string { return TableNamePainter }
