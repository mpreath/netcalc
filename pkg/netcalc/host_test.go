package netcalc

import (
	"encoding/json"
	"testing"
)

func TestMarshalJSONHost(t *testing.T) {
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
		testHost, err := NewHost(testAddress, testMask)
		if err != nil {
			t.Fatalf(err.Error())
		}

		s, err := json.Marshal(testHost)
		if err != nil {
			t.Errorf(err.Error())
		}

		if len(s) <= 0 {
			t.Errorf("didn't receive any output from marshal")
		}
	}
}

func TestNewHost(t *testing.T) {
	testCases := []struct {
		ddAddress     string
		ddMask        string
		errorExpected bool
	}{
		{"192.168.1.1", "255.255.255.0", false},
		{"192.168.1.1", "255.255.255.134", true},
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
		testHost, err := NewHost(testAddress, testMask)
		if err != nil {
			if !testCase.errorExpected {
				t.Fatalf(err.Error())
			} else {
				continue
			}
		}

		ddTestAddress := ExportAddress(testHost.Address)
		ddTestMask := ExportAddress(testHost.Mask)

		if ddTestAddress != testCase.ddAddress {
			t.Errorf("generated address (%s) doesn't match spec address (%s)", ddTestAddress, testCase.ddAddress)
		}

		if ddTestMask != testCase.ddMask {
			t.Errorf("generated mask (%s) doesn't match spec mask (%s)", ddTestMask, testCase.ddMask)
		}
	}
}
