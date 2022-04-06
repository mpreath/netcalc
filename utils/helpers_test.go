package utils

import "testing"

func TestIsValidMask(t *testing.T) {
	test_cases := []struct {
		dd_mask  string
		is_valid bool
	}{
		{"255.255.255.0", true},
		{"192.168.1.0", false},
		{"128.0.0.0", true},
		{"0.0.0.128", false},
	}

	for _, test_case := range test_cases {
		uint_mask, _ := Ddtoi(test_case.dd_mask)
		is_mask_valid := IsValidMask(uint_mask)
		if is_mask_valid != test_case.is_valid {
			t.Errorf("result for %s (%t) does not match spec (%t)", test_case.dd_mask, is_mask_valid, test_case.is_valid)
		}

	}
}
