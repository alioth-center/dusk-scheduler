package entity

type NoRequestContent struct{}

type BaseResponse[T any] struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}
