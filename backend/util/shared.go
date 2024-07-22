package util 

func Contains(arr []string, str string) bool {
	
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func ContainsInt64(arr []int64, num int64) bool {
	for _, a := range arr {
		if a == num {
			return true
		}
	}
	return false
}
