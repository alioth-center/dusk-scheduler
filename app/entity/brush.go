package entity

import "encoding/json"

type RegisterBrushRequest struct {
	Maintainer string `json:"maintainer" binding:"required"`
	Protocol   string `json:"protocol" binding:"required"`
}

type RegisterBrushResponse struct {
	Name   string              `json:"name"`
	Secret string              `json:"secret"`
	Policy RegisterBrushPolicy `json:"policy"`
}

type RegisterBrushPolicy struct {
	Protocol string          `json:"protocol"`
	Options  json.RawMessage `json:"options"`
}
