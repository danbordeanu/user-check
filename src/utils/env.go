package utils

import (
	"os"
	"strconv"
)

func EnvOrDefault(name, def string) string {
	if v, ok := os.LookupEnv(name); ok && v != "" {
		return v
	}
	return def
}

func EnvOrDefaultInt32(name string, def int32) int32 {
	if v, ok := os.LookupEnv(name); ok && v != "" {
		vc, _ := strconv.ParseInt(v, 10, 32)
		return int32(vc)
	}
	return def
}



