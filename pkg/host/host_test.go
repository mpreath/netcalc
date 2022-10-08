package host

import (
	"encoding/json"
	"testing"

	"github.com/mpreath/netcalc/pkg/utils"
)

func TestMarshalJSON(t *testing.T) {
	test_cases := []struct {
		dd_address string
		dd_mask    string
	}{
		{"192.168.1.1", "255.255.255.0"},
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
		test_host, err := New(testAddress, testMask)
		if err != nil {
			t.Fatalf(err.Error())
		}

		s, err := json.Marshal(test_host)
		if err != nil {
			t.Errorf(err.Error())
		}

		if len(s) <= 0 {
			t.Errorf("didn't receive any output from marshal")
		}
	}
}

func TestNew(t *testing.T) {
	test_cases := []struct {
		dd_address string
		dd_mask    string
	}{
		{"192.168.1.1", "255.255.255.0"},
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
		test_host, err := New(testAddress, testMask)
		if err != nil {
			t.Fatalf(err.Error())
		}

		dd_test_address := utils.Itodd(test_host.Address)
		dd_test_mask := utils.Itodd(test_host.Mask)

		if dd_test_address != test_case.dd_address {
			t.Errorf("generated address (%s) doesn't match spec address (%s)", dd_test_address, test_case.dd_address)
		}

		if dd_test_mask != test_case.dd_mask {
			t.Errorf("generated mask (%s) doesn't match spec mask (%s)", dd_test_mask, test_case.dd_mask)
		}
	}
}
