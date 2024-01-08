package utils

import "strconv"

func ValidateEAN(ean string) bool {
	length := len(ean)

	if length != 8 && length != 13 {
		return false
	}

	sum := 0
	for i := 0; i < length-1; i++ {
		num, err := strconv.Atoi(string(ean[i]))
		if err != nil {
			return false
		}

		if (i+length)%2 == 0 {
			sum += num * 3
		} else {
			sum += num
		}
	}

	checkDigit, err := strconv.Atoi(string(ean[length-1]))
	if err != nil {
		return false
	}

	return (10-sum%10)%10 == checkDigit
}
