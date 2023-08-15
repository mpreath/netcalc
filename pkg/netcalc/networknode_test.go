package netcalc

import (
	"testing"
)

func TestSplit(t *testing.T) {
	testCases := []struct {
		ddAddress   string
		ddMask      string
		errorString string
	}{
		{"192.168.1.1", "255.255.255.0", ""},
		{"10.1.0.0", "255.255.0.0", ""},
		{"10.1.1.1", "255.255.255.252", "NetworkNode#Split: network doesn't support being split"},
	}

	for _, testCase := range testCases {
		testAddress, err := ParseAddress(testCase.ddAddress)
		if err != nil {
			t.Fatalf(err.Error())
		}
		testMask, err := ParseAddress(testCase.ddMask)
		if err != nil {
			t.Fatalf(err.Error())
		}

		testNetwork, _ := NewNetwork(testAddress, testMask)
		testNode := NewNetworkNode(testNetwork)

		err = testNode.Split()
		if err != nil {
			if err.Error() != testCase.errorString {
				t.Errorf(err.Error())
			}
			continue
		}

		if len(testNode.Subnets) == 0 {
			t.Errorf("unable to split network (%s %s)", testCase.ddAddress, testCase.ddMask)
		}
	}
}

func TestSplitToHostCount(t *testing.T) {
	testCases := []struct {
		ddAddress         string
		ddMask            string
		hostCount         int
		expectedHostCount int
		errorString       string
	}{
		{"192.168.1.1", "255.255.255.0", 1, 2, ""},
		{"10.1.0.0", "255.255.0.0", 100, 126, ""},
		{"192.168.1.1", "255.255.255.0", 300, 0, "ValidForHostCount: network can't support that many hosts"},
	}

	for _, testCase := range testCases {
		testNetworkAddress, err := ParseAddress(testCase.ddAddress)
		if err != nil {
			t.Fatalf(err.Error())
		}
		testNetworkMask, err := ParseAddress(testCase.ddMask)
		if err != nil {
			t.Fatalf(err.Error())
		}
		testNetwork, _ := NewNetwork(testNetworkAddress, testNetworkMask)
		testNode := NewNetworkNode(testNetwork)

		err = SplitToHostCount(testNode, testCase.hostCount)

		if err != nil {
			if err.Error() != testCase.errorString {
				t.Errorf(err.Error())
			}
			continue
		}

		// traverse to a leaf networknode
		for len(testNode.Subnets) > 0 {
			testNode = testNode.Subnets[0]
		}

		if testNode.Network.HostCount() != testCase.expectedHostCount {
			t.Errorf("subnet host count (%d) doesn't match spec (%d)", testNode.Network.HostCount(), testCase.hostCount)
		}
	}
}

func TestSplitToNetCount(t *testing.T) {
	testCases := []struct {
		ddAddress        string
		ddMask           string
		netCount         int
		expectedNetCount int
		errorString      string
	}{
		{"192.168.1.1", "255.255.255.0", 2, 2, ""},
		{"10.1.0.0", "255.255.0.0", 4, 4, ""},
		{"192.168.1.1", "255.255.255.252", 2, 0, "SplitToNetCount: network can't support that many subnetworks"},
	}

	for _, testCase := range testCases {
		testNetworkAddress, err := ParseAddress(testCase.ddAddress)
		if err != nil {
			t.Fatalf(err.Error())
		}
		testNetworkMask, err := ParseAddress(testCase.ddMask)
		if err != nil {
			t.Fatalf(err.Error())
		}
		testNetwork, _ := NewNetwork(testNetworkAddress, testNetworkMask)
		testNode := NewNetworkNode(testNetwork)

		err = SplitToNetCount(testNode, testCase.netCount)

		if err != nil {
			if err.Error() != testCase.errorString {
				t.Errorf(err.Error())
			}
			continue
		}

		netCount := testNode.NetworkCount()
		if netCount != testCase.expectedNetCount {
			t.Errorf("network count (%d) doesn't match spec (%d)", netCount, testCase.expectedNetCount)
		}
	}
}

func TestSplitToVlsmCount(t *testing.T) {
	testCases := []struct {
		ddAddress     string
		ddMask        string
		vlsmList      []int
		errorExpected bool
	}{
		{"192.168.1.1", "255.255.255.0", []int{2, 2, 2, 2}, false},
		{"10.1.0.0", "255.255.0.0", []int{2, 2, 2, 2}, false},
		{"192.168.1.1", "255.255.255.252", []int{2, 2, 2, 2}, true},
	}

	for _, testCase := range testCases {
		testNetworkAddress, err := ParseAddress(testCase.ddAddress)
		if err != nil {
			t.Fatalf(err.Error())
		}
		testNetworkMask, err := ParseAddress(testCase.ddMask)
		if err != nil {
			t.Fatalf(err.Error())
		}
		testNetwork, _ := NewNetwork(testNetworkAddress, testNetworkMask)
		testNode := NewNetworkNode(testNetwork)

		for _, vlsm := range testCase.vlsmList {
			err = SplitToVlsmCount(testNode, vlsm)
			if err != nil {
				if !testCase.errorExpected {
					t.Fatalf(err.Error())
				}
				return
			}
		}

		// flatten to just get network list
		networkList := testNode.FlattenUtilized()

		if len(networkList) != len(testCase.vlsmList) {
			t.Errorf("network list length (%d) !=  vlsm list length(%d)", len(networkList), len(testCase.vlsmList))
		}
	}
}
