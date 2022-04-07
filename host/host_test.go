package host

import (
	"encoding/json"
	"testing"

	"github.com/mpreath/netcalc/utils"
)

func TestMarshalJSON(t *testing.T) {
	test_cases := []struct {
		dd_address string
		dd_mask    string
	}{
		{"192.168.1.1", "255.255.255.0"},
	}

	for _, test_case := range test_cases {
		test_host, _ := GenerateHost(test_case.dd_address, test_case.dd_mask)

		s, err := json.Marshal(test_host)
		if err != nil {
			t.Errorf(err.Error())
		}

		if len(s) <= 0 {
			t.Errorf("didn't receive any output from marshal")
		}
	}
}

func TestGenerateHost(t *testing.T) {
	test_cases := []struct {
		dd_address string
		dd_mask    string
	}{
		{"192.168.1.1", "255.255.255.0"},
	}

	for _, test_case := range test_cases {
		test_host, _ := GenerateHost(test_case.dd_address, test_case.dd_mask)

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
