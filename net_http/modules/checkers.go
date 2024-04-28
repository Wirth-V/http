package modules

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	re = regexp.MustCompile("[^a-zA-Z0-9]+")
)

func check(s ...string) error {
	var sb strings.Builder
	var control int = Zero
	for id := range s {
		//проверка на допустимые символы
		res := re.FindAllString(s[id], -1)

		if res != nil {
			control++
			sb.WriteString("in the argument ")
			sb.WriteString(s[id])
			sb.WriteString(" invalid characters are present\t")
		}
		//проверка длинны
		if len(s[id]) >= Size {
			control++
			sb.WriteString("the argument ")
			sb.WriteString(s[id])
			sb.WriteString(" exceeds the maximum size of the input data (no more than 20 characters\t)")
		}
	}

	if control != Zero {
		return fmt.Errorf(sb.String())
	}
	return nil
}
