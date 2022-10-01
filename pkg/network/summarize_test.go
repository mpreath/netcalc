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
		networks, _ := SplitToHostCount(test_network, 2)
		summarized_network, _ := SummarizeNetworks(networks)

		if summarized_network.Address != test_network.Address {
			t.Errorf("summarized network doesn't match test network")
		}

		if summarized_network.Mask != test_network.Mask {
			t.Errorf("summarized mask doesn't match test mask")
		}

	}

}
