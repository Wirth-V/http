package moduls

import (
	"regexp"
)

func Sanitize(s string) bool {
	if len(s) >= 21 {
		ErrorLog.Println("Превышен максимальный размер вхожных данных (не более 8 символов)")
		return true
	}
	if b, _ := regexp.MatchString("[^a-zA-Z0-9]+", s); b {
		ErrorLog.Println("Присутсвуют недопустимые символы")
		return true
	}
	return false
}
