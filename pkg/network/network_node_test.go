package network

import (
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
		test_network, _ := GenerateNetwork(test_case.dd_address, test_case.dd_mask)
		test_node := &NetworkNode{
			Network: test_network,
		}

		err := test_node.Split()

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
