package domain

const TableNameClient = "client"

type Client struct{}

func (Client) TableName() string { return TableNameClient }
