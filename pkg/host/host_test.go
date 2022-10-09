package host

import (
	"encoding/json"
	"testing"

	"github.com/mpreath/netcalc/pkg/utils"
)

func TestMarshalJSON(t *testing.T) {
	testCases := []struct {
		ddAddress string
		ddMask    string
	}{
		{"192.168.1.1", "255.255.255.0"},
	}

	for _, testCase := range testCases {
		testAddress, err := utils.ParseAddress(testCase.ddAddress)
		if err != nil {
			t.Fatalf(err.Error())
		}

		testMask, err := utils.ParseAddress(testCase.ddMask)
		if err != nil {
			t.Fatalf(err.Error())
		}
		testHost, err := New(testAddress, testMask)
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

func TestNew(t *testing.T) {
	testCases := []struct {
		ddAddress string
		ddMask    string
	}{
		{"192.168.1.1", "255.255.255.0"},
	}

	for _, testCase := range testCases {
		testAddress, err := utils.ParseAddress(testCase.ddAddress)
		if err != nil {
			t.Fatalf(err.Error())
		}

		testMask, err := utils.ParseAddress(testCase.ddMask)
		if err != nil {
			t.Fatalf(err.Error())
		}
		testHost, err := New(testAddress, testMask)
		if err != nil {
			t.Fatalf(err.Error())
		}

		ddTestAddress := utils.ExportAddress(testHost.Address)
		ddTestMask := utils.ExportAddress(testHost.Mask)

		if ddTestAddress != testCase.ddAddress {
			t.Errorf("generated address (%s) doesn't match spec address (%s)", ddTestAddress, testCase.ddAddress)
		}

		if ddTestMask != testCase.ddMask {
			t.Errorf("generated mask (%s) doesn't match spec mask (%s)", ddTestMask, testCase.ddMask)
		}
	}
}
