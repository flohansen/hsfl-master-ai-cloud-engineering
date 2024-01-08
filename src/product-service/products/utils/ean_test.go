package utils

import "testing"

func TestValidateEAN(t *testing.T) {
	testCases := []struct {
		name    string
		ean     string
		isValid bool
	}{
		{"Valid EAN8", "96385074", true},
		{"Invalid EAN8", "12345675", false},
		{"Valid EAN13", "4006381333931", true},
		{"Invalid EAN13", "4006381333932", false},
		{"Invalid Length", "123", false},
		{"Non Numeric", "abcdefg", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateEAN(tc.ean)
			if result != tc.isValid {
				t.Errorf("For EAN '%s', expected: %v, got: %v", tc.ean, tc.isValid, result)
			}
		})
	}
}
