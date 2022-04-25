package basefunc

import (
	"regexp"
)

func CheckIdCard(idCardNum string) bool {
	res, err := regexp.MatchString("(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)", idCardNum)
	if err != nil {
		return false
	}
	return res
}
