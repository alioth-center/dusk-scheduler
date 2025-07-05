package entity

type RegisterPainterRequest struct {
	Maintainer string `json:"maintainer" binding:"required"`
	Slot       int    `json:"slot" binding:"required"`
	Name       string `json:"name,omitempty" binding:"omitempty"`
	Abilities  string `json:"abilities,omitempty" binding:"omitempty"`
	Protocol   string `json:"protocol,omitempty" binding:"required_with=Abilities"`
}

type RegisterPainterResponse struct {
	Name   string                      `json:"name"`
	Secret string                      `json:"secret"`
	Policy RegisterPainterUploadPolicy `json:"policy"`
}

type RegisterPainterUploadPolicy struct {
	Protocol string                    `json:"protocol"`
	S3       *RegisterPainterS3Policy  `json:"s3,omitempty"`
	Ftp      *RegisterPainterFtpPolicy `json:"ftp,omitempty"`
}

type RegisterPainterS3Policy struct {
	Protocol  string `json:"protocol"`
	Endpoint  string `json:"endpoint"`
	Bucket    string `json:"bucket"`
	Prefix    string `json:"prefix"`
	Region    string `json:"region"`
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
}

type RegisterPainterFtpPolicy struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	RemotePath string `json:"remote_path"`
}

type DeletePainterRequest = NoRequestContent

type DeletePainterResponse = NoRequestContent

type GetPainterTaskListRequest = NoRequestContent

type GetPainterTaskListResponse struct {
	Tasks []GetPainterTaskListRequest `json:"tasks"`
}

type GetPainterTaskListItem struct {
	TaskID         int    `json:"task_id"`
	Height         int    `json:"height"`
	Width          int    `json:"width"`
	Priority       int    `json:"priority"`
	DelaySeconds   int    `json:"delay_seconds"`
	EncodedContent string `json:"encoded_content"`
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
