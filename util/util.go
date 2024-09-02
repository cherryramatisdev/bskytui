package util

import "os"

func IsDebug() bool {
	return os.Getenv("DEBUG") == "1"
}
