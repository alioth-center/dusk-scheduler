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

type GetPainterTaskListResponse struct{}

type CompletePainterTaskRequest struct{}

type CompletePainterTaskResponse struct{}
