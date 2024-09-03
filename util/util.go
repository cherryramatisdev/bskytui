package util

import "os"

func IsDebug() bool {
	return os.Getenv("DEBUG") == "1"
}

func GetCountLabel(count int, label string) string {

	if count == 1 {
		return label
	} else {
		if label == "reply" {
			return "replies"
		}
		return label + "s"
	}
}
