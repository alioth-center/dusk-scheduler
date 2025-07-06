package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"github.com/alioth-center/dusk-scheduler/app/service/errors"
	"io"
)

func EncryptHmacSha256String(message, secret string) string {
	key, data := []byte(secret), []byte(message)

	h := hmac.New(sha256.New, key)
	h.Write(data)
	digest := h.Sum(nil)

	return base64.StdEncoding.EncodeToString(digest)
}

func EncryptMd5String(message string) string {
	data := []byte(message)

	h := md5.New()
	h.Write(data)
	digest := h.Sum(nil)

	return base64.StdEncoding.EncodeToString(digest)
}

func EncryptBase64String(message string) string {
	return base64.StdEncoding.EncodeToString([]byte(message))
}

func UnifyEncodingToBase64(message string, encoding string) (result string, err error) {
	switch encoding {
	case "base64":
		return message, nil
	case "gzip":
		reader, decodeErr := gzip.NewReader(bytes.NewBufferString(message))
		if decodeErr != nil {
			return "", decodeErr
		}

		defer errors.Ignore(reader.Close())
		rawBytes, readErr := io.ReadAll(reader)
		if readErr != nil {
			return "", readErr
		}

		return base64.StdEncoding.EncodeToString(rawBytes), nil
	case "raw":
		return base64.StdEncoding.EncodeToString([]byte(message)), nil
	default:
		return "", errors.CustomError("unknown encoding: " + encoding)
	}
}
