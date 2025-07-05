package entity

type NoRequestContent struct{}

type NoResponseContent struct{}

type BaseResponse[T any] struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

func SuccessResponse[T any](data T) BaseResponse[T] {
	return BaseResponse[T]{
		Status: "success",
		Code:   200,
		Data:   data,
	}
}
