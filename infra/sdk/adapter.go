package sdk

import (
	"bytes"
	"context"
	"net/http"
)

type BrushSDK interface {
	RenderImage(ctx context.Context, protocol string, callURL string, encodedContent *bytes.Buffer) (output *bytes.Buffer, err error)
}

func NewBrushSDK(_ *http.Client) BrushSDK {
	return nil
}
