package entity

type RegisterRequest struct {
	EmailAddress   string `json:"email_address" binding:"required"`
	RedemptionCode string `json:"redemption_code,omitempty" binding:"omitempty"`
}

type RegisterResponse struct {
	ExpiredAt int64 `json:"expired_at"`
}

type AuthorizeRequest struct {
	EmailAddress      string `json:"email_address" binding:"required"`
	AuthorizationCode string `json:"authorization_code" binding:"required"`
}

type AuthorizeResponse struct {
	Maintainer string `json:"maintainer"`
	ApiKey     string `json:"api_key"`
}

type GetMetadataRequest = NoRequestContent

type GetMetadataResponse struct {
	Maintainer string                  `json:"maintainer"`
	ApiKey     string                  `json:"api_key"`
	Options    GetMetadataClientOption `json:"options"`
}

type GetMetadataClientOption struct {
	BrushApiEnable          bool `json:"brush_api_enable"`
	DelayRenderEnable       bool `json:"delay_render_enable"`
	MaxRenderHeight         int  `json:"max_render_height"`
	MaxRenderWidth          int  `json:"max_render_width"`
	MaxPriority             int  `json:"max_priority"`
	MaxRenderSize           int  `json:"max_render_size"`
	RequestFrequency        int  `json:"request_frequency"`
	FrequencyIntervalSecond int  `json:"frequency_interval_second"`
}

type GetCompletedTasksRequest = NoRequestContent

type GetCompletedTasksResponse struct {
	Tasks   []GetCompletedTaskItem `json:"tasks"`
	HasMore bool                   `json:"has_more,omitempty"`
}

type GetCompletedTaskItem struct {
	TaskID        int                     `json:"task_id"`
	Size          string                  `json:"size"`
	Priority      int                     `json:"priority"`
	ContentHash   string                  `json:"content_hash"`
	Status        string                  `json:"status"`
	Timestamps    GetTaskStatusTimestamps `json:"timestamps"`
	ArchiveReason string                  `json:"archive_reason,omitempty"`
}

type GetCompletedTaskItemTimestamps struct {
	CreatedAt   int64 `json:"created_at"`
	ScheduledAt int64 `json:"scheduled_at,omitempty"`
	CompletedAt int64 `json:"completed_at,omitempty"`
	ArchivedAt  int64 `json:"archived_at,omitempty"`
}

type GetQuotaRequest = NoRequestContent

type GetQuotaResponse struct {
	Details        GetQuotaDetails `json:"details"`
	LastCheckpoint int64           `json:"last_checkpoint"`
}

type GetQuotaDetails struct {
	TotalQuota     int `json:"total_quota"`
	UsedQuota      int `json:"used_quota"`
	RemainingQuota int `json:"remaining_quota"`
}
