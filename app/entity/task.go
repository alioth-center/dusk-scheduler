package entity

type CreateTaskRequest struct {
	Content         string `json:"content" binding:"required"`
	ContentEncoding string `json:"content_encoding" binding:"required"`
	RenderWidth     int    `json:"render_width" binding:"required"`
	RenderHeight    int    `json:"render_height" binding:"required"`
	OutputFormat    string `json:"output_format,omitempty" binding:"omitempty"`
	DelaySeconds    int    `json:"delay_seconds,omitempty" binding:"omitempty"`
}

type CreateTaskResponse = NoResponseContent

type GetTaskStatusRequest = NoRequestContent

type GetTaskStatusResponse struct {
	Size             string                   `json:"size"`
	Priority         string                   `json:"priority"`
	ContentHash      string                   `json:"content_hash"`
	Status           string                   `json:"status"`
	Timestamps       *GetTaskStatusTimestamps `json:"timestamps"`
	OutcomeReference string                   `json:"outcome_reference,omitempty"`
	ArchiveReason    string                   `json:"archive_reason,omitempty"`
}

type GetTaskStatusTimestamps struct {
	CreatedAt   int64 `json:"created_at"`
	ScheduledAt int64 `json:"scheduled_at,omitempty"`
	CompletedAt int64 `json:"completed_at,omitempty"`
	ArchivedAt  int64 `json:"archived_at,omitempty"`
}
