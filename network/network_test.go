package network

import (
	"testing"

	"github.com/mpreath/netcalc/utils"
)

func TestGenerateNetwork(t *testing.T) {
	test_cases := []struct {
		dd_address         string
		dd_mask            string
		dd_network_address string
	}{
		{"192.168.1.1", "255.255.255.0", "192.168.1.0"},
	}

	for _, test_case := range test_cases {
		test_network, _ := GenerateNetwork(test_case.dd_address, test_case.dd_mask)

		dd_test_network := utils.Itodd(test_network.Address)
		dd_test_mask := utils.Itodd(test_network.Mask)

		if dd_test_network != test_case.dd_network_address {
			t.Errorf("generated address (%s) doesn't match spec address (%s)", dd_test_network, test_case.dd_network_address)
		}

		if dd_test_mask != test_case.dd_mask {
			t.Errorf("generated mask (%s) doesn't match spec mask (%s)", dd_test_mask, test_case.dd_mask)
		}
	}
}

func TestGenerateNetworkFromBits(t *testing.T) {
	test_cases := []struct {
		dd_network_address string
		dd_mask            string
	}{
		{"192.168.1.0", "255.255.255.0"},
	}

	for _, test_case := range test_cases {
		tmp_network, _ := GenerateNetwork(test_case.dd_network_address, test_case.dd_mask)

		test_network, _ := GenerateNetworkFromBits(tmp_network.Address, tmp_network.Mask)

		if test_network.Address != tmp_network.Address {
			t.Errorf("generated address (%s) doesn't match spec address (%s)", utils.Itodd(test_network.Address), utils.Itodd(tmp_network.Address))
		}

		if test_network.Mask != tmp_network.Mask {
			t.Errorf("generated mask (%s) doesn't match spec mask (%s)", utils.Itodd(test_network.Mask), utils.Itodd(tmp_network.Mask))
		}
	}
}
