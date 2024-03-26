package moduls

func Sanitize(s string) {
	if len(s) >= 9 {
		ErrorLog.Println("Превышен максимальный размер вхожных данных (не более 8 символов)")
	}

}
