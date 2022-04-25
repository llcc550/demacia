package basefunc

import "strings"

func RemoveSpace(in string) string {
	return strings.Replace(strings.Replace(in, "\r\n", "", -1), " ", "", -1)
}
