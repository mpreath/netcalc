package utils

import "testing"

func TestDDtoi(t *testing.T) {
	test_cases := []struct {
		dd_address   string
		uint_address uint32
		error_string string
	}{
		{"255.255.255.255", 4294967295, ""},
		{"0.0.0.0", 0, ""},
		{"192.168.1.1", 3232235777, ""},
		{"192.168.1.256", 0, "utils:ParseAddress: parsing \"256\": number must be 255 or less"},
		{"192.168.1.A", 0, "strconv.ParseUint: parsing \"A\": invalid syntax"},
		{"192.168.1.-1", 0, "strconv.ParseUint: parsing \"-1\": invalid syntax"},
		{"192.168.1.1.8", 0, "utils:ParseAddress: parsing \"192.168.1.1.8\": too many octets"},
	}

	for _, test_case := range test_cases {
		test_val, err := ParseAddress(test_case.dd_address)
		if err != nil {
			// error encountered
			if err.Error() != test_case.error_string {
				t.Fatalf("result (%s) does not match spec (%s)", err.Error(), test_case.error_string)
			}
		}

		if test_val != test_case.uint_address {
			// function calculation is incorrect
			t.Errorf("result (%d) does not match spec (%d)", test_val, test_case.uint_address)
		}
	}

}

func TestItodd(t *testing.T) {
	test_cases := []struct {
		dd_address   string
		uint_address uint32
		error_string string
	}{
		{"255.255.255.255", 4294967295, ""},
		{"0.0.0.0", 0, ""},
		{"192.168.1.1", 3232235777, ""},
	}

	for _, test_case := range test_cases {
		test_val := Itodd(test_case.uint_address)

		if test_val != test_case.dd_address {
			// function calculation is incorrect
			t.Errorf("result (%s) does not match spec (%s)", test_val, test_case.dd_address)
		}
	}
}
