package network

import (
	"encoding/json"
	"github.com/mpreath/netcalc/pkg/host"
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
		test_network, _ := host.New(testAddress, testMask)

		s, err := json.Marshal(test_network)
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
		dd_address         string
		dd_mask            string
		dd_network_address string
		error_string       string
	}{
		{"192.168.1.1", "255.255.255.0", "192.168.1.0", ""},
		{"192.168.1.1", "255.0.255.0", "192.168.1.0", "network.GenerateNetwork: invalid subnet mask"},
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
		test_network, err := host.New(testAddress, testMask)
		if err != nil {
			if err.Error() != test_case.error_string {
				t.Errorf(err.Error())
			}
			continue
		}

		dd_test_network := utils.Itodd(test_network.Address)
		dd_test_mask := utils.Itodd(test_network.Mask)

		if dd_test_network != test_case.dd_network_address {
			t.Errorf("generated address (%s) doesn't match spec address (%s)", dd_test_network, test_case.dd_network_address)
		}

		if dd_test_mask != test_case.dd_mask {
			t.Errorf("generated mask (%s) doesn't match spec mask (%s)", dd_test_mask, test_case.dd_mask)
		}
	}
}

func TestGetHosts(t *testing.T) {
	test_cases := []struct {
		dd_address string
		dd_mask    string
		host_count int
	}{
		{"192.168.1.0", "255.255.255.0", 254},
		{"192.168.1.0", "255.255.255.128", 126},
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
		test_network, err := New(testAddress, testMask)
		if err != nil {
			t.Fatalf(err.Error())
		}

		test_hosts := test_network.Hosts()

		if len(test_hosts) != test_case.host_count {
			t.Errorf("generated host count (%d) doesn't match spec count (%d)", len(test_hosts), test_case.host_count)
		}
	}
}
