package network

import (
	"testing"
)

func TestSummarizeNetworks(t *testing.T) {
	type NetworkMap struct {
		Address string
		Mask    string
	}
	testCases := []struct {
		testNetworks    []NetworkMap
		expectedNetwork NetworkMap
	}{
		{[]NetworkMap{{Address: "192.168.1.0", Mask: "255.255.255.252"},
			{Address: "192.168.1.4", Mask: "255.255.255.252"}},
			NetworkMap{Address: "192.168.1.0", Mask: "255.255.255.248"}},
	}

	for _, testCase := range testCases {
		var testNetworks []*Network
		expectedNetwork, err := New(testCase.expectedNetwork.Address, testCase.expectedNetwork.Mask)
		if err != nil {
			t.Fatalf(err.Error())
		}
		for _, networkMap := range testCase.testNetworks {
			testNetwork, err := New(networkMap.Address, networkMap.Mask)
			if err != nil {
				t.Fatalf(err.Error())
			}
			testNetworks = append(testNetworks, testNetwork)
		}

		summarizedNetwork, err := SummarizeNetworks(testNetworks)

		if summarizedNetwork.Address != expectedNetwork.Address {
			t.Errorf("summarized network doesn't match test network")
		}

		if summarizedNetwork.Mask != expectedNetwork.Mask {
			t.Errorf("summarized mask doesn't match test mask")
		}

	}

}
