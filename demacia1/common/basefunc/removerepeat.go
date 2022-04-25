package basefunc

func RemoveRepByLoop(slc []int64) []int64 {
	result := make([]int64, 0, len(slc))
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false
				break
			}
		}
		if flag {
			result = append(result, slc[i])
		}
	}
	return result
}
