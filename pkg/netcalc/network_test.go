package netcalc

import (
	"encoding/json"
	"testing"
)

func TestMarshalJSONNetwork(t *testing.T) {
	testCases := []struct {
		ddAddress string
		ddMask    string
	}{
		{"192.168.1.1", "255.255.255.0"},
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

		s, err := json.Marshal(testNetwork)
		if err != nil {
			t.Errorf(err.Error())
		}

		if len(s) <= 0 {
			t.Errorf("didn't receive any output from marshal")
		}
	}
}

func TestUnmarshalJSONNetwork(t *testing.T) {
	testCases := []struct {
		ddAddress  string
		jsonString string
	}{
		{"192.168.1.0", "{ \"address\": \"192.168.1.0\", \"mask\": \"255.255.255.0\" }"},
	}

	for _, testCase := range testCases {
		var testNetwork Network
		err := json.Unmarshal([]byte(testCase.jsonString), &testNetwork)
		if err != nil {
			t.Errorf(err.Error())
		}

		expectedResult, err := ParseAddress(testCase.ddAddress)
		if err != nil {
			t.Errorf(err.Error())
		}

		if testNetwork.Address != expectedResult {
			t.Errorf("unmarshalled address (%s) doesn't match spec address (%s)", ExportAddress(testNetwork.Address), ExportAddress(expectedResult))

		}
	}
}

func TestNewNetwork(t *testing.T) {
	testCases := []struct {
		ddAddress        string
		ddMask           string
		ddNetworkAddress string
		errorString      string
	}{
		{"192.168.1.1", "255.255.255.0", "192.168.1.0", ""},
		{"192.168.1.1", "255.0.255.0", "192.168.1.0", "network.New: invalid subnet mask"},
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
		testNetwork, err := NewNetwork(testAddress, testMask)
		if err != nil {
			if err.Error() != testCase.errorString {
				t.Errorf(err.Error())
			}
			continue
		}

		ddTestNetwork := ExportAddress(testNetwork.Address)
		ddTestMask := ExportAddress(testNetwork.Mask)

		if ddTestNetwork != testCase.ddNetworkAddress {
			t.Errorf("generated address (%s) doesn't match spec address (%s)", ddTestNetwork, testCase.ddNetworkAddress)
		}

		if ddTestMask != testCase.ddMask {
			t.Errorf("generated mask (%s) doesn't match spec mask (%s)", ddTestMask, testCase.ddMask)
		}
	}
}

func TestHosts(t *testing.T) {
	testCases := []struct {
		ddAddress string
		ddMask    string
		hostCount int
	}{
		{"192.168.1.0", "255.255.255.0", 254},
		{"192.168.1.0", "255.255.255.128", 126},
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
		testNetwork, err := NewNetwork(testAddress, testMask)
		if err != nil {
			t.Fatalf(err.Error())
		}

		testHosts := testNetwork.Hosts()

		if len(testHosts) != testCase.hostCount {
			t.Errorf("generated host count (%d) doesn't match spec count (%d)", len(testHosts), testCase.hostCount)
		}
	}
}

func TestHostCount(t *testing.T) {
	testCases := []struct {
		ddAddress string
		ddMask    string
		hostCount int
	}{
		{"192.168.1.0", "255.255.255.0", 254},
		{"192.168.1.0", "255.255.255.128", 126},
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
		testNetwork, err := NewNetwork(testAddress, testMask)
		if err != nil {
			t.Fatalf(err.Error())
		}

		testHostCount := testNetwork.HostCount()

		if testHostCount != testCase.hostCount {
			t.Errorf("generated host count (%d) doesn't match spec count (%d)", testHostCount, testCase.hostCount)
		}
	}
}

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
		expectedNetworkAddress, err := ParseAddress(testCase.expectedNetwork.Address)
		if err != nil {
			t.Fatalf(err.Error())
		}
		expectedNetworkMask, err := ParseAddress(testCase.expectedNetwork.Mask)
		if err != nil {
			t.Fatalf(err.Error())
		}
		expectedNetwork, err := NewNetwork(expectedNetworkAddress, expectedNetworkMask)
		if err != nil {
			t.Fatalf(err.Error())
		}
		for _, networkMap := range testCase.testNetworks {
			testNetworkAddress, err := ParseAddress(networkMap.Address)
			if err != nil {
				t.Fatalf(err.Error())
			}
			testNetworkMask, err := ParseAddress(networkMap.Mask)
			if err != nil {
				t.Fatalf(err.Error())
			}
			testNetwork, err := NewNetwork(testNetworkAddress, testNetworkMask)
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
