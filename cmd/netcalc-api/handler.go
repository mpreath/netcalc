package main

import (
	"encoding/json"
	"github.com/mpreath/netcalc/pkg/network"
	"github.com/mpreath/netcalc/pkg/network/networknode"
	"github.com/mpreath/netcalc/pkg/utils"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Response struct {
	Status    string      `json:"status"`
	Error     string      `json:"error,omitempty"`
	ErrorCode int         `json:"error_code,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func Info(w http.ResponseWriter, r *http.Request) {

	ipAddress, err := utils.ParseAddress(r.URL.Query().Get("address"))
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	subnetMask, err := utils.ParseAddress(r.URL.Query().Get("mask"))
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	net, err := network.New(ipAddress, subnetMask)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	response := Response{Status: "ok", Data: net}
	writeJsonResponse(w, http.StatusOK, response)

}

func Subnet(w http.ResponseWriter, r *http.Request) {

	ipAddress, err := utils.ParseAddress(r.URL.Query().Get("address"))
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	subnetMask, err := utils.ParseAddress(r.URL.Query().Get("mask"))
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	hostCount, _ := strconv.Atoi(r.URL.Query().Get("hostCount"))
	networkCount, _ := strconv.Atoi(r.URL.Query().Get("networkCount"))

	net, err := network.New(ipAddress, subnetMask)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	node := networknode.New(net)

	if hostCount > 0 {
		err = networknode.SplitToHostCount(node, hostCount)
		if err != nil {
			log.Fatal(err)
		}

	} else if networkCount > 0 {
		err = networknode.SplitToNetCount(node, networkCount)
		if err != nil {
			log.Fatal(err)
		}
	}

	response := Response{Status: "ok", Data: flattenResults(node)}
	writeJsonResponse(w, http.StatusOK, response)
}

func Summarize(w http.ResponseWriter, r *http.Request) {
	var networkList []*network.Network

	err := json.NewDecoder(r.Body).Decode(&networkList)
	if err != nil {
		writeErrorResponse(w, err)
	}
	summarizedNetwork, err := network.SummarizeNetworks(networkList)
	if err != nil {
		writeErrorResponse(w, err)
	}
	response := Response{Status: "ok", Data: summarizedNetwork}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		writeErrorResponse(w, err)
	}
}

func Vlsm(w http.ResponseWriter, r *http.Request) {

	ipAddress, err := utils.ParseAddress(r.URL.Query().Get("address"))
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	subnetMask, err := utils.ParseAddress(r.URL.Query().Get("mask"))
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	net, err := network.New(ipAddress, subnetMask)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
	// generate network from args
	node := &networknode.NetworkNode{
		Network: net,
	}

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
		err = networknode.SplitToVlsmCount(node, vlsm)

		if err != nil {
			writeErrorResponse(w, err)
			return
		}
	}

	response := Response{Status: "ok", Data: flattenResults(node)}
	writeJsonResponse(w, http.StatusOK, response)
}

func flattenResults(node *networknode.NetworkNode) []*network.Network {
	var networkList []*network.Network
	if len(node.Subnets) == 0 {
		return append(networkList, node.Network)
	} else {
		networkList = append(networkList, flattenResults(node.Subnets[0])...)
		networkList = append(networkList, flattenResults(node.Subnets[1])...)
		return networkList
	}
}

func writeJsonResponse(w http.ResponseWriter, status int, response Response) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(jsonResponse)
}

func writeErrorResponse(w http.ResponseWriter, err error) {
	response := Response{}
	response.Status = "error"
	response.Error = err.Error()
	response.ErrorCode = http.StatusInternalServerError
	writeJsonResponse(w, response.ErrorCode, response)
}
