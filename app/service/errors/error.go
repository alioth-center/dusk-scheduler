package errors

import (
	"fmt"
)

func Ignore(_ error) {}

type ApiError struct{}

func (e ApiError) Error() string {
	return ""
}

func BadRequestError(bindErr error) ApiError {
	return ApiError{}
}

func InvalidParameter(parameterName string) error {
	return fmt.Errorf("invalid parameter: %s", parameterName)
}

func ResourceNotFound(resource string, reference any) error {
	return fmt.Errorf("resource %s not found: %v", resource, reference)
}

func InternalError() ApiError {
	return ApiError{}
}
