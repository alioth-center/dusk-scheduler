package entity

import "encoding/json"

type RegisterPainterRequest struct {
	Maintainer string `json:"maintainer" binding:"required"`
	Slot       int    `json:"slot" binding:"required"`
}

type RegisterPainterResponse struct {
	Name   string                      `json:"name"`
	Secret string                      `json:"secret"`
	Policy RegisterPainterUploadPolicy `json:"policy"`
}

type RegisterPainterUploadPolicy struct {
	Protocol string          `json:"protocol"`
	Options  json.RawMessage `json:"options"`
}

type ReconnectPainterRequest = NoRequestContent

type ReconnectPainterResponse = NoResponseContent

type DeletePainterRequest = NoRequestContent

type DeletePainterResponse = NoRequestContent

type GetPainterTaskListRequest = NoRequestContent

type GetPainterTaskListResponse struct {
	Tasks []GetPainterTaskListItem `json:"tasks"`
}

type GetPainterTaskListItem struct {
	TaskID         int    `json:"task_id"`
	Height         int    `json:"height"`
	Width          int    `json:"width"`
	Priority       int    `json:"priority"`
	DelaySeconds   int    `json:"delay_seconds"`
	EncodedContent string `json:"encoded_content"`
	Checksum       string `json:"checksum"`
}

type CompletePainterTaskRequest struct {
	Status           string `json:"status" binding:"required"`
	Message          string `json:"message" binding:"required"`
	StorageReference string `json:"storage_reference" binding:"required_if=Status complete"`
	StartedAt        int64  `json:"started_at" binding:"required"`
	CompletedAt      int64  `json:"completed_at" binding:"required"`
}

type CompletePainterTaskResponse struct {
	Next []GetPainterTaskListItem `json:"next,omitempty"`
}
