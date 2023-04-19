package fullname

import (
	"fmt"
)

func Convert(first string, last string) *string {
	var str string
	str = fmt.Sprintf("%s %s", first, last)
	return &str
}
