package moduls

import (
	"regexp"
)

func Sanitize(s string) bool {
	if b, _ := regexp.MatchString("[^a-zA-Z0-9]+", s); b {
		ErrorLog.Println("Присутсвуют недопустимые символы")
		return true
	}
	return false
}

func Length(s string) bool {
	if len(s) >= 21 {
		ErrorLog.Println("Превышен максимальный размер вхожных данных (не более 20 символов)")
		return true
	}
	return false
}
