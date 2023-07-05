package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/mpreath/netcalc/pkg/netcalc"
)

type Response struct {
	Status    string      `json:"status"`
	Error     string      `json:"error,omitempty"`
	ErrorCode int         `json:"error_code,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type NetworkInfo struct {
	NetworkAddress   string `json:"network_address"`
	Mask             string `json:"mask"`
	BroadcastAddress string `json:"broadcast_address"`
	NumberOfHosts    int    `json:"number_of_hosts"`
}

func Info(w http.ResponseWriter, r *http.Request) {

	ipAddress, err := netcalc.ParseAddress(r.URL.Query().Get("address"))
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	subnetMask, err := netcalc.ParseAddress(r.URL.Query().Get("mask"))
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	net, err := netcalc.NewNetwork(ipAddress, subnetMask)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	info := &NetworkInfo{
		NetworkAddress:   netcalc.ExportAddress(net.Address),
		Mask:             netcalc.ExportAddress(net.Mask),
		BroadcastAddress: netcalc.ExportAddress(net.BroadcastAddress()),
		NumberOfHosts:    net.HostCount(),
	}

	writeJsonResponse(w, http.StatusOK, info)

}

func Subnet(w http.ResponseWriter, r *http.Request) {

	ipAddress, err := netcalc.ParseAddress(r.URL.Query().Get("address"))
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	subnetMask, err := netcalc.ParseAddress(r.URL.Query().Get("mask"))
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	hostCount, _ := strconv.Atoi(r.URL.Query().Get("hosts"))
	networkCount, _ := strconv.Atoi(r.URL.Query().Get("networks"))

	if hostCount == 0 && networkCount == 0 {
		writeErrorResponse(w, fmt.Errorf("subnet: no host or network counts provided"))
		return
	}

	net, err := netcalc.NewNetwork(ipAddress, subnetMask)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	node := netcalc.NewNetworkNode(net)

	if hostCount > 0 {
		err = netcalc.SplitToHostCount(node, hostCount)
		if err != nil {
			log.Fatal(err)
		}

	} else if networkCount > 0 {
		err = netcalc.SplitToNetCount(node, networkCount)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		writeErrorResponse(w, fmt.Errorf("no valid hosts or networks value provided"))
		return
	}

	writeJsonResponse(w, http.StatusOK, node.Flatten())
}

func Summarize(w http.ResponseWriter, r *http.Request) {
	var networkList []*netcalc.Network

	err := json.NewDecoder(r.Body).Decode(&networkList)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
	summarizedNetwork, err := netcalc.SummarizeNetworks(networkList)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	writeJsonResponse(w, http.StatusOK, summarizedNetwork)
}

func Vlsm(w http.ResponseWriter, r *http.Request) {

	ipAddress, err := netcalc.ParseAddress(r.URL.Query().Get("address"))
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	subnetMask, err := netcalc.ParseAddress(r.URL.Query().Get("mask"))
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	net, err := netcalc.NewNetwork(ipAddress, subnetMask)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	node := netcalc.NewNetworkNode(net)

	vlsmArgs := strings.Split(r.URL.Query().Get("vlsmList"), ",")
	var vlsmList = make([]int, len(vlsmArgs))
	for idx, val := range vlsmArgs {
		vlsmList[idx], err = strconv.Atoi(val)
		if err != nil {
			writeErrorResponse(w, err)
			return
		}
	}
	sort.Slice(vlsmList, func(i, j int) bool {
		return vlsmList[i] < vlsmList[j]
	})

	for _, vlsm := range vlsmList {
		err = netcalc.SplitToVlsmCount(node, vlsm)

		if err != nil {
			writeErrorResponse(w, err)
			return
		}
	}

	writeJsonResponse(w, http.StatusOK, node.Flatten())
}

func writeJsonResponse(w http.ResponseWriter, status int, data interface{}) {
	enc := writeHeaders(w, status)
	err := enc.Encode(Response{
		Status: "ok",
		Data:   data,
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func writeErrorResponse(w http.ResponseWriter, err error) {
	enc := writeHeaders(w, http.StatusInternalServerError)
	err = enc.Encode(Response{
		Status:    "error",
		Error:     err.Error(),
		ErrorCode: http.StatusInternalServerError,
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func writeHeaders(w http.ResponseWriter, status int) *json.Encoder {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	return enc
}
