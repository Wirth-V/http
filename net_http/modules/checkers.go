package modules

import (
	"fmt"
	"regexp"
	"strings"
)

func check(s ...string) error {
	var sb strings.Builder
	var control int = Zero
	for id := range s {
		//проверка на допустимые символы
		re, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			return fmt.Errorf("the argument %v could not be checked for the validity of characters", s[id])
		}
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
