package networknode

import (
	"github.com/mpreath/netcalc/pkg/network"
	"github.com/mpreath/netcalc/pkg/utils"
	"testing"
)

func TestSplit(t *testing.T) {
	test_cases := []struct {
		dd_address   string
		dd_mask      string
		error_string string
	}{
		{"192.168.1.1", "255.255.255.0", ""},
		{"10.1.0.0", "255.255.0.0", ""},
		{"10.1.1.1", "255.255.255.252", "network:Split: network doesn't support being split"},
	}

	for _, test_case := range test_cases {
		testAddress, err := utils.ParseAddress(test_case.dd_address)
		if err != nil {
			t.Fatalf(err.Error())
		}
		testMask, err := utils.ParseAddress(test_case.dd_mask)
		if err != nil {
			t.Fatalf(err.Error())
		}

		test_network, _ := network.New(testAddress, testMask)
		test_node := &NetworkNode{
			Network: test_network,
		}

		err = test_node.Split()
		if err != nil {
			if err.Error() != test_case.error_string {
				t.Errorf(err.Error())
			}
			continue
		}

		if len(test_node.Subnets) == 0 {
			t.Errorf("unable to split network (%s %s)", test_case.dd_address, test_case.dd_mask)
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
		testNetworkAddress, err := utils.ParseAddress(test_case.dd_address)
		if err != nil {
			t.Fatalf(err.Error())
		}
		testNetworkMask, err := utils.ParseAddress(test_case.dd_mask)
		if err != nil {
			t.Fatalf(err.Error())
		}
		test_network, _ := network.New(testNetworkAddress, testNetworkMask)
		test_node := New(test_network)

		err = SplitToHostCount(test_node, test_case.host_count)

		if err != nil {
			if err.Error() != test_case.error_string {
				t.Errorf(err.Error())
			}
			continue
		}

		// traverse to a leaf networknode
		for len(test_node.Subnets) > 0 {
			test_node = test_node.Subnets[0]
		}

		if test_node.Network.HostCount() != test_case.expected_host_count {
			t.Errorf("subnet host count (%d) doesn't match spec (%d)", test_node.Network.HostCount(), test_case.host_count)
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
		testNetworkAddress, err := utils.ParseAddress(test_case.dd_address)
		if err != nil {
			t.Fatalf(err.Error())
		}
		testNetworkMask, err := utils.ParseAddress(test_case.dd_mask)
		if err != nil {
			t.Fatalf(err.Error())
		}
		test_network, _ := network.New(testNetworkAddress, testNetworkMask)
		test_node := New(test_network)

		err = SplitToNetCount(test_node, test_case.net_count)

		if err != nil {
			if err.Error() != test_case.error_string {
				t.Errorf(err.Error())
			}
			continue
		}

		net_count := test_node.NetworkCount()
		if net_count != test_case.expected_net_count {
			t.Errorf("network count (%d) doesn't match spec (%d)", net_count, test_case.expected_net_count)
		}
	}
}
