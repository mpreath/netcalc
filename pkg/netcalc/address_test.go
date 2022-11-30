package netcalc

import (
	"testing"
)

func TestIsValidMask(t *testing.T) {
	testCases := []struct {
		ddMask  string
		isValid bool
	}{
		{"255.255.255.0", true},
		{"192.168.1.0", false},
		{"128.0.0.0", true},
		{"0.0.0.128", false},
	}

	for _, testCase := range testCases {
		uintMask, _ := ParseAddress(testCase.ddMask)
		isMaskValid := IsValidMask(uintMask)
		if isMaskValid != testCase.isValid {
			t.Errorf("result for %s (%t) does not match spec (%t)", testCase.ddMask, isMaskValid, testCase.isValid)
		}

	}
}

func TestGetNetworkAddress(t *testing.T) {
	testCases := []struct {
		ddAddress        string
		ddMask           string
		ddNetworkAddress string
	}{
		{"192.168.1.1", "255.255.255.0", "192.168.1.0"},
	}

	for _, testCase := range testCases {
		address, _ := ParseAddress(testCase.ddAddress)
		mask, _ := ParseAddress(testCase.ddMask)
		networkAddress, _ := ParseAddress(testCase.ddNetworkAddress)

		testNetworkAddress := GetNetworkAddress(address, mask)

		if testNetworkAddress != networkAddress {
			t.Errorf("network %s doesn't match spec network %s", ExportAddress(testNetworkAddress), ExportAddress(networkAddress))
		}

	}
}

func TestGetBroadcastAddress(t *testing.T) {
	testCases := []struct {
		ddAddress          string
		ddMask             string
		ddBroadcastAddress string
	}{
		{"192.168.1.1", "255.255.255.0", "192.168.1.255"},
	}

	for _, testCase := range testCases {
		address, _ := ParseAddress(testCase.ddAddress)
		mask, _ := ParseAddress(testCase.ddMask)
		broadcastAddress, _ := ParseAddress(testCase.ddBroadcastAddress)

		testBroadcastAddress := GetBroadcastAddress(address, mask)

		if testBroadcastAddress != broadcastAddress {
			t.Errorf("broadcast %s doesn't match spec broadcast %s", ExportAddress(testBroadcastAddress), ExportAddress(broadcastAddress))
		}

	}
}

func TestGetBitsInMask(t *testing.T) {
	testCases := []struct {
		ddMask       string
		expectedBits int
	}{
		{"255.255.255.252", 30},
		{"255.255.255.255", 32},
		{"128.0.0.0", 1},
	}

	for _, testCase := range testCases {
		mask, _ := ParseAddress(testCase.ddMask)

		testBitsInMask := GetBitsInMask(mask)

		if testBitsInMask != testCase.expectedBits {
			t.Errorf("mask %s [%d] doesn't contain the expected number of bits [%d]", testCase.ddMask, testBitsInMask, testCase.expectedBits)
		}

	}
}

func TestGetCommonBitMask(t *testing.T) {
	testCases := []struct {
		testNetworks          [2]string
		expectedCommonBitMask string
	}{
		{[2]string{"192.168.1.0", "192.168.2.0"}, "255.255.252.0"},
	}

	for _, testCase := range testCases {
		testNetwork1, err := ParseAddress(testCase.testNetworks[0])
		if err != nil {
			t.Fatalf(err.Error())
		}
		testNetwork2, err := ParseAddress(testCase.testNetworks[1])
		if err != nil {
			t.Fatalf(err.Error())
		}

		commonBitMask := GetCommonBitMask(testNetwork1, testNetwork2)
		expectedBitMask, err := ParseAddress(testCase.expectedCommonBitMask)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if commonBitMask != expectedBitMask {
			t.Errorf("results (%s) don't match expectations (%s)", ExportAddress(commonBitMask), ExportAddress(expectedBitMask))
		}

	}
}

func TestParseAddress(t *testing.T) {
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

func TestExportAddress(t *testing.T) {
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
