package entity

type GetOutcomeRequest = NoRequestContent

type GetOutcomeResponse struct {
	Metadata GetOutcomeMetadata `json:"metadata"`
	Content  GetOutcomeContent  `json:"content"`
}

type GetOutcomeMetadata struct {
	TaskReference int    `json:"task_reference"`
	InstanceName  string `json:"instance_name"`
	StartedAt     int64  `json:"started_at"`
	CompletedAt   int64  `json:"completed_at"`
}

type GetOutcomeContent struct {
	Base64Encoded string `json:"base64_encoded,omitempty"`
	ImageURL      string `json:"image_url,omitempty"`
}

type AcknowledgeRequest = NoRequestContent

type AcknowledgeResponse = NoResponseContent
