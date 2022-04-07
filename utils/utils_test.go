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
		uint_mask, _ := Ddtoi(test_case.dd_mask)
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
		address, _ := Ddtoi(test_case.dd_address)
		mask, _ := Ddtoi(test_case.dd_mask)
		network_address, _ := Ddtoi(test_case.dd_network_address)

		test_network_address := GetNetworkAddress(address, mask)

		if test_network_address != network_address {
			t.Errorf("network %s doesn't match spec network %s", Itodd(test_network_address), Itodd(network_address))
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
		address, _ := Ddtoi(test_case.dd_address)
		mask, _ := Ddtoi(test_case.dd_mask)
		broadcast_address, _ := Ddtoi(test_case.dd_broadcast_address)

		test_broadcast_address := GetBroadcastAddress(address, mask)

		if test_broadcast_address != broadcast_address {
			t.Errorf("broadcast %s doesn't match spec broadcast %s", Itodd(test_broadcast_address), Itodd(broadcast_address))
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
		mask, _ := Ddtoi(test_case.dd_mask)

		test_bits_in_mask := GetBitsInMask(mask)

		if test_bits_in_mask != test_case.expected_bits {
			t.Errorf("mask %s [%d] doesn't contain the expected number of bits [%d]", test_case.dd_mask, test_bits_in_mask, test_case.expected_bits)
		}

	}
}

func TestGetMaskFromBits(t *testing.T) {
	test_cases := []struct {
		dd_mask      string
		bits         int
		error_string string
	}{
		{"255.255.255.252", 30, ""},
		{"255.255.255.255", 32, ""},
		{"128.0.0.0", 1, ""},
		{"0.0.0.0", 33, "utils:GetMaskFromBits: bits must be 32 or less"},
	}

	for _, test_case := range test_cases {

		test_case_mask, _ := Ddtoi(test_case.dd_mask)
		mask, err := GetMaskFromBits(test_case.bits)

		if err != nil {
			// error encountered
			if err.Error() != test_case.error_string {
				t.Fatalf("result (%s) does not match spec (%s)", err.Error(), test_case.error_string)
			}
			continue
		}

		if mask != test_case_mask {
			t.Errorf("generated mask (%s) doesn't match expected mask (%s)", Itodd(mask), test_case.dd_mask)
		}

	}
}
