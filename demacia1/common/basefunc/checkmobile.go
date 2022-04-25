package basefunc

import "regexp"

func CheckMobile(in string) bool {
	mobileIsOk, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, in)
	return mobileIsOk
}
