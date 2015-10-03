package utils

import (
	"net/http"
	"strings"
)

func GetIpFromRequest(request *http.Request) string {
	remoteAddress := request.RemoteAddr
	header := request.Header
	xRealIp := header.Get("X-Real-Ip")
	xForwardedFor := header.Get("X-Forwarded-For")
	switch {
	case xForwardedFor != "":
		parts := strings.Split(xForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// TODO: should return first non-local address
		return parts[0]
	case xRealIp != "":
		return xRealIp
	default:
		index := strings.LastIndex(remoteAddress, ":")
		if index == -1 {
			return remoteAddress
		}
		return remoteAddress[:index]
	}
}
