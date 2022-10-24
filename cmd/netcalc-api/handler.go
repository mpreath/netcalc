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

	writeJsonResponse(w, http.StatusOK, net)

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

	hostCount, err := strconv.Atoi(r.URL.Query().Get("hostCount"))
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	networkCount, err := strconv.Atoi(r.URL.Query().Get("networkCount"))
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

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

	writeJsonResponse(w, http.StatusOK, node.Flatten())
}

func Summarize(w http.ResponseWriter, r *http.Request) {
	var networkList []*network.Network

	err := json.NewDecoder(r.Body).Decode(&networkList)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
	summarizedNetwork, err := network.SummarizeNetworks(networkList)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	writeJsonResponse(w, http.StatusOK, summarizedNetwork)
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

	node := networknode.New(net)

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

	writeJsonResponse(w, http.StatusOK, node.Flatten())
}

func writeJsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
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
