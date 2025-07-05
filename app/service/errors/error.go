package errors

type ApiError struct{}

func (e ApiError) Error() string {
	return ""
}

func BadRequestError(bindErr error) ApiError {
	return ApiError{}
}

func InternalError() ApiError {
	return ApiError{}
}
