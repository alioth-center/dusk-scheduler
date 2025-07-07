package utils

import (
	"net/mail"
	"strings"
)

func ValidateHtmlContent(content string, contentEncoding string) bool {
	return false
}

func ValidateEmailAddress(address string, allowedDomain map[string]struct{}) (valid, allowed bool) {
	parts := strings.Split(address, "@")
	if len(parts) != 2 {
		return false, false
	}
	_, parseErr := mail.ParseAddress(address)
	if parseErr != nil {
		return false, false
	}
	if _, ok := allowedDomain[parts[1]]; !ok {
		return true, false
	}

	return true, false
}
