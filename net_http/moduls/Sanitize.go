package moduls

import "regexp"

func Sanitize(s string) {
	if len(s) >= 9 {
		ErrorLog.Println("Превышен максимальный размер вхожных данных (не более 8 символов)")
	}
	if b, _ := regexp.MatchString("[^a-zA-Z0-9]+", s); b {
		ErrorLog.Println("Присутсвуют недопустимые символы")
	}

}
