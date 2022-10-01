package network

import (
	"encoding/json"
	"testing"

	"github.com/mpreath/netcalc/pkg/utils"
)

func TestMarshalJSON(t *testing.T) {
	test_cases := []struct {
		dd_address string
		dd_mask    string
	}{
		{"192.168.1.1", "255.255.255.0"},
	}

	for _, test_case := range test_cases {
		test_network, _ := GenerateNetwork(test_case.dd_address, test_case.dd_mask)

		s, err := json.Marshal(test_network)
		if err != nil {
			t.Errorf(err.Error())
		}

		if len(s) <= 0 {
			t.Errorf("didn't receive any output from marshal")
		}
	}
}

func TestGenerateNetwork(t *testing.T) {
	test_cases := []struct {
		dd_address         string
		dd_mask            string
		dd_network_address string
		error_string       string
	}{
		{"192.168.1.1", "255.255.255.0", "192.168.1.0", ""},
		{"192.168.1.1", "255.0.255.0", "192.168.1.0", "network.GenerateNetwork: invalid subnet mask"},
	}

	for _, test_case := range test_cases {
		test_network, err := GenerateNetwork(test_case.dd_address, test_case.dd_mask)
		if err != nil {
			if err.Error() != test_case.error_string {
				t.Errorf(err.Error())
			}
			continue
		}

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

func TestGetHosts(t *testing.T) {
	test_cases := []struct {
		dd_network_address string
		dd_mask            string
		host_count         int
	}{
		{"192.168.1.0", "255.255.255.0", 254},
		{"192.168.1.0", "255.255.255.128", 126},
	}

	for _, test_case := range test_cases {
		test_network, _ := GenerateNetwork(test_case.dd_network_address, test_case.dd_mask)
		test_hosts := GetHosts(test_network)

		if len(test_hosts) != test_case.host_count {
			t.Errorf("generated host count (%d) doesn't match spec count (%d)", len(test_hosts), test_case.host_count)
		}
	}
}

func TestSplitToHostCount(t *testing.T) {
	test_cases := []struct {
		dd_address          string
		dd_mask             string
		host_count          int
		expected_host_count int
		error_string        string
	}{
		{"192.168.1.1", "255.255.255.0", 1, 2, ""},
		{"10.1.0.0", "255.255.0.0", 100, 126, ""},
		{"192.168.1.1", "255.255.255.0", 300, 0, "network.SplitToHostCount: network can't support that many hosts"},
	}

	for _, test_case := range test_cases {
		test_network, _ := GenerateNetwork(test_case.dd_address, test_case.dd_mask)

		subnets, err := SplitToHostCount(test_network, test_case.host_count)

		if err != nil {
			if err.Error() != test_case.error_string {
				t.Errorf(err.Error())
			}
			continue
		}

		if subnets[0].HostCount() != test_case.expected_host_count {
			t.Errorf("subnet host count (%d) doesn't match spec (%d)", subnets[0].HostCount(), test_case.host_count)
		}
	}
}

func TestSplitToNetCount(t *testing.T) {
	test_cases := []struct {
		dd_address         string
		dd_mask            string
		net_count          int
		expected_net_count int
		error_string       string
	}{
		{"192.168.1.1", "255.255.255.0", 2, 2, ""},
		{"10.1.0.0", "255.255.0.0", 4, 4, ""},
		{"192.168.1.1", "255.255.255.252", 2, 0, "network.SplitToNetCount: network can't support that many subnetworks"},
	}

	for _, test_case := range test_cases {
		test_network, _ := GenerateNetwork(test_case.dd_address, test_case.dd_mask)

		subnets, err := SplitToNetCount(test_network, test_case.net_count)

		if err != nil {
			if err.Error() != test_case.error_string {
				t.Errorf(err.Error())
			}
			continue
		}

		net_count := len(subnets)
		if net_count != test_case.expected_net_count {
			t.Errorf("network count (%d) doesn't match spec (%d)", net_count, test_case.expected_net_count)
		}
	}
}
