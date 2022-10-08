package network

import (
	"github.com/mpreath/netcalc/pkg/host"
	"github.com/mpreath/netcalc/pkg/utils"
	"testing"
)

func TestSummarizeNetworks(t *testing.T) {
	type NetworkMap struct {
		Address string
		Mask    string
	}
	testCases := []struct {
		testNetworks    []NetworkMap // networks to be summarized
		expectedNetwork NetworkMap   // what the summarization should be
	}{
		{
			[]NetworkMap{
				{Address: "192.168.1.0", Mask: "255.255.255.252"},
				{Address: "192.168.1.4", Mask: "255.255.255.252"},
			},
			NetworkMap{
				Address: "192.168.1.0", Mask: "255.255.255.248",
			},
		},
		{
			[]NetworkMap{
				{Address: "192.168.1.0", Mask: "255.255.255.192"},
				{Address: "192.168.1.64", Mask: "255.255.255.192"},
				{Address: "192.168.1.128", Mask: "255.255.255.192"},
				{Address: "192.168.1.192", Mask: "255.255.255.192"},
			},
			NetworkMap{
				Address: "192.168.1.0", Mask: "255.255.255.0",
			},
		},
	}

	for _, testCase := range testCases {
		var testNetworks []*Network
		expectedNetworkAddress, err := utils.ParseAddress(testCase.expectedNetwork.Address)
		if err != nil {
			t.Fatalf(err.Error())
		}
		expectedNetworkMask, err := utils.ParseAddress(testCase.expectedNetwork.Mask)
		if err != nil {
			t.Fatalf(err.Error())
		}
		expectedNetwork, err := host.New(expectedNetworkAddress, expectedNetworkMask)
		if err != nil {
			t.Fatalf(err.Error())
		}
		for _, networkMap := range testCase.testNetworks {
			testNetworkAddress, err := utils.ParseAddress(networkMap.Address)
			if err != nil {
				t.Fatalf(err.Error())
			}
			testNetworkMask, err := utils.ParseAddress(networkMap.Mask)
			if err != nil {
				t.Fatalf(err.Error())
			}
			testNetwork, err := New(testNetworkAddress, testNetworkMask)
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
