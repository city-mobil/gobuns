package config

import "strings"

func parseConfigType(filePath string) string {
	sp := strings.Split(filePath, ".")
	if len(sp) == 0 {
		return ""
	}
	return sp[len(sp)-1]
}

func SanitizePrefix(p string) string {
	if n := len(p); n > 0 && p[n-1] != '.' {
		return p + "."
	}
	return p
}
