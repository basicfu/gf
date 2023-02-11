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
		str := fw
		if i := strings.IndexByte(str, ','); i >= 0 {
			str = str[:i]
		}
		if i := strings.LastIndexByte(str, ':'); i >= 0 {
			str = str[i+1:] //截取本地ip情况::ffff:10.0.0.2
		}
		return str
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

func HasLocalIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip.IsLoopback() {
		return true
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}
	return ip4[0] == 10 || // 10.0.0.0/8
		(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || // 172.16.0.0/12
		(ip4[0] == 169 && ip4[1] == 254) || // 169.254.0.0/16
		(ip4[0] == 192 && ip4[1] == 168) // 192.168.0.0/16
}
