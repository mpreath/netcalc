package utils

import (
	"testing"
)

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
		uint_mask, _ := ParseAddress(test_case.dd_mask)
		is_mask_valid := IsValidMask(uint_mask)
		if is_mask_valid != test_case.is_valid {
			t.Errorf("result for %s (%t) does not match spec (%t)", test_case.dd_mask, is_mask_valid, test_case.is_valid)
		}

	}
}

func TestGetNetworkAddress(t *testing.T) {
	test_cases := []struct {
		dd_address         string
		dd_mask            string
		dd_network_address string
	}{
		{"192.168.1.1", "255.255.255.0", "192.168.1.0"},
	}

	for _, test_case := range test_cases {
		address, _ := ParseAddress(test_case.dd_address)
		mask, _ := ParseAddress(test_case.dd_mask)
		network_address, _ := ParseAddress(test_case.dd_network_address)

		test_network_address := GetNetworkAddress(address, mask)

		if test_network_address != network_address {
			t.Errorf("network %s doesn't match spec network %s", ExportAddress(test_network_address), ExportAddress(network_address))
		}

	}
}

func TestGetBroadcastAddress(t *testing.T) {
	test_cases := []struct {
		dd_address           string
		dd_mask              string
		dd_broadcast_address string
	}{
		{"192.168.1.1", "255.255.255.0", "192.168.1.255"},
	}

	for _, test_case := range test_cases {
		address, _ := ParseAddress(test_case.dd_address)
		mask, _ := ParseAddress(test_case.dd_mask)
		broadcast_address, _ := ParseAddress(test_case.dd_broadcast_address)

		test_broadcast_address := GetBroadcastAddress(address, mask)

		if test_broadcast_address != broadcast_address {
			t.Errorf("broadcast %s doesn't match spec broadcast %s", ExportAddress(test_broadcast_address), ExportAddress(broadcast_address))
		}

	}
}

func TestGetBitsInMask(t *testing.T) {
	test_cases := []struct {
		dd_mask       string
		expected_bits int
	}{
		{"255.255.255.252", 30},
		{"255.255.255.255", 32},
		{"128.0.0.0", 1},
	}

	for _, test_case := range test_cases {
		mask, _ := ParseAddress(test_case.dd_mask)

		test_bits_in_mask := GetBitsInMask(mask)

		if test_bits_in_mask != test_case.expected_bits {
			t.Errorf("mask %s [%d] doesn't contain the expected number of bits [%d]", test_case.dd_mask, test_bits_in_mask, test_case.expected_bits)
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
