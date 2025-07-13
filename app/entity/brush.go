package entity

import "encoding/json"

type RegisterBrushRequest struct {
	Maintainer string `json:"maintainer" binding:"required"`
	Protocol   string `json:"protocol" binding:"required"`
	CallURL    string `json:"call_url" binding:"required"`
}

type RegisterBrushResponse struct {
	BrushID int `json:"brush_id"`
}

type RegisterBrushPolicy struct {
	Protocol string          `json:"protocol"`
	Options  json.RawMessage `json:"options"`
}
