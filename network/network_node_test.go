package network

import (
	"testing"
)

func TestSplitToHostCount(t *testing.T) {
	test_cases := []struct {
		dd_address          string
		dd_mask             string
		host_count          int
		expected_host_count int
	}{
		{"192.168.1.1", "255.255.255.0", 1, 2},
		{"10.1.0.0", "255.255.0.0", 100, 126},
	}

	for _, test_case := range test_cases {
		test_network, _ := GenerateNetwork(test_case.dd_address, test_case.dd_mask)
		test_node := &NetworkNode{
			Network: test_network,
		}

		SplitToHostCount(test_node, test_case.host_count)

		// traverse to a leaf node
		for len(test_node.Subnets) > 0 {
			test_node = test_node.Subnets[0]
		}

		if len(test_node.Network.Hosts) != test_case.expected_host_count {
			t.Errorf("subnet host count (%d) doesn't match spec (%d)", len(test_node.Network.Hosts), test_case.host_count)
		}
	}
}

func TestSplitToNetCount(t *testing.T) {
	test_cases := []struct {
		dd_address         string
		dd_mask            string
		net_count          int
		expected_net_count int
	}{
		{"192.168.1.1", "255.255.255.0", 2, 2},
		{"10.1.0.0", "255.255.0.0", 4, 4},
	}

	for _, test_case := range test_cases {
		test_network, _ := GenerateNetwork(test_case.dd_address, test_case.dd_mask)
		test_node := &NetworkNode{
			Network: test_network,
		}

		SplitToNetCount(test_node, test_case.net_count)

		net_count := GetNetworkCount(test_node)
		if net_count != uint(test_case.expected_net_count) {
			t.Errorf("network count (%d) doesn't match spec (%d)", net_count, test_case.expected_net_count)
		}
	}
}
