package netcalc

import "testing"

func TestDDtoi(t *testing.T) {
	testCases := []struct {
		ddAddress   string
		uintAddress uint32
		errorString string
	}{
		{"255.255.255.255", 4294967295, ""},
		{"0.0.0.0", 0, ""},
		{"192.168.1.1", 3232235777, ""},
		{"192.168.1.256", 0, "utils:ParseAddress: parsing \"256\": number must be 255 or less"},
		{"192.168.1.A", 0, "strconv.ParseUint: parsing \"A\": invalid syntax"},
		{"192.168.1.-1", 0, "strconv.ParseUint: parsing \"-1\": invalid syntax"},
		{"192.168.1.1.8", 0, "utils:ParseAddress: parsing \"192.168.1.1.8\": too many octets"},
	}

	for _, testCase := range testCases {
		testVal, err := ParseAddress(testCase.ddAddress)
		if err != nil {
			// error encountered
			if err.Error() != testCase.errorString {
				t.Fatalf("result (%s) does not match spec (%s)", err.Error(), testCase.errorString)
			}
		}

		if testVal != testCase.uintAddress {
			// function calculation is incorrect
			t.Errorf("result (%d) does not match spec (%d)", testVal, testCase.uintAddress)
		}
	}

}

func TestItodd(t *testing.T) {
	testCases := []struct {
		ddAddress   string
		uintAddress uint32
		errorString string
	}{
		{"255.255.255.255", 4294967295, ""},
		{"0.0.0.0", 0, ""},
		{"192.168.1.1", 3232235777, ""},
	}

	for _, testCase := range testCases {
		testVal := ExportAddress(testCase.uintAddress)

		if testVal != testCase.ddAddress {
			// function calculation is incorrect
			t.Errorf("result (%s) does not match spec (%s)", testVal, testCase.ddAddress)
		}
	}
}
