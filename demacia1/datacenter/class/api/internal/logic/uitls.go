package logic

import (
	"errors"
	"strconv"
)

func IntToCa(n int) (string, error) {
	nStr := strconv.Itoa(n)
	var company = []string{"", "十", "百", "千", "万"}
	var zhCa = []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}

	var res string
	if len(nStr) > len(company) {
		return "", errors.New("the length of parameter 1 is out of range")
	}
	zero := false
	for i := 1; i <= len(nStr); i++ {
		site, _ := strconv.Atoi(nStr[i-1 : i])
		if i == len(nStr) && site == 0 {
			break
		}
		if site == 0 {
			if !zero {
				res += zhCa[site]
				zero = true
			}
		} else {
			res += zhCa[site] + company[len(nStr)-i]
			zero = false
		}
	}
	return res, nil
}
