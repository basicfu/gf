package gip

import (
	"net"
	"net/http"
	"strings"
)

func RealIP(r *http.Request) string {
	if contextIp := r.Context().Value("remote_addr"); contextIp != nil {
		return contextIp.(string)
	}
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	if fw := r.Header.Get("X-Forwarded-For"); fw != "" {
		if i := strings.IndexByte(fw, ','); i >= 0 {
			return fw[:i]
		}
		return fw
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}
