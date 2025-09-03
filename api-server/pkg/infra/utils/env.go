package utils

import "os"

//LONGHORN_UPLOAD_URL_PREFIX

func GetEnv(variable string, defaultValue *string) string {
	value, exists := os.LookupEnv(variable)
	if !exists {
		return *defaultValue
	}
	return value
}
