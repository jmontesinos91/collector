package middleware

import (
	"net/http"
	"strings"

	"github.com/jmontesinos91/osecurity/sts"
)

type Paths string

const (
	read         Paths = "/v1/traffic"
	export       Paths = "/v1/traffic/export"
	resetcounter Paths = "/v1/traffic/counter"
)

func ValidatePermission(permission sts.Permission, path string, method string) bool {
	action := strings.ReplaceAll(permission.Action, "_", "")
	action = strings.ToLower(action)

	switch action {
	case "read":
		if strings.Contains(string(read), path) && method == http.MethodGet {
			return true
		}
	case "export":
		if strings.Contains(string(export), path) && method == http.MethodGet {
			return true
		}
	case "resetcounter":
		if strings.Contains(string(resetcounter), path) && method == http.MethodPost {
			return true
		}
	default:
		return false
	}

	return false
}
