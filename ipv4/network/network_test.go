package network

import "testing"

func TestGenerateNetwork(t *testing.T) {
	test_cases := []struct {
		dd_address   string
		uint_address uint32
		dd_mask      string
		uint_mask    uint32
	}{
		{"192.168.1.0", 3232235776, "255.255.255.0", 4294967040},
	}

	for _, test_case := range test_cases {
		test_network, _ := GenerateNetwork(test_case.dd_address, test_case.dd_mask)

		if test_network.Address != test_case.uint_address {
			t.Errorf("generated address (%d) doesn't match spec address (%d)", test_network.Address, test_case.uint_address)
		}

		if test_network.Mask != test_case.uint_mask {
			t.Errorf("generated mask (%d) doesn't match spec mask (%d)", test_network.Mask, test_case.uint_mask)
		}
	}
}
