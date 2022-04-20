package network

import (
	"testing"
)

func TestSummarizeNetworks(t *testing.T) {
	test_cases := []struct {
		dd_address string
		mask       string
	}{
		{"192.168.1.0", "255.255.255.0"},
		{"10.1.0.0", "255.255.0.0"},
		{"172.16.12.0", "255.255.248.0"},
	}

	for _, test_case := range test_cases {
		test_network, _ := GenerateNetwork(test_case.dd_address, test_case.mask)
		test_network_node := &NetworkNode{
			Network: test_network,
		}
		SplitToHostCount(test_network_node, 2)
		networks := NetworkNodeToArray(test_network_node)
		summarized_networks := SummarizeNetworks(networks)

		if len(summarized_networks) != 1 {
			t.Errorf("incorrect number of networks (%d) in summary", len(summarized_networks))
		}

	}

}
