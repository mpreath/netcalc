![Netcalc Logo](https://github.com/mpreath/public_files/blob/25a6d05318b3ddfbf7b3d3835201aa1a0cfa48ab/netcalc.png)

## Overview
netcalc is an advanced CLI-based network subnetting tool. It provides the
following features:

* provided an ip address and mask, it will calculate the network address, broadcast address, bits in mask, and all hosts of the network
* provided an ip address, mask, and network count, it will calculate the network addresses, and masks in order to subnet the network into smaller networks based on the number of networks desired
* provided an ip address, mask, and host count, it will calculate the network addresses, and mask in ord to subnet the network into smaller networks based on the number of hosts desired 
* provided an ip address, mask, and a list of host count values, it will calculate a list of VLSM subnets that match the list of host count requirements. Instead of making the subnets all the same size, it only creates enough subnets to match the list, and they may all be different sizes
* provided a list of subnets, it will try to summarize the subnets into one (or more) supernets if possible.

*NOTE:* netcalc uses bit shifting and mask operators to work with ip addresses in their actual 32 bit integer form. It converts from dotted decimal notation to unsigned integer values, and then manipulates the integer values.

*NOTE:* netcalc uses a binary tree structure to subnet networks by extending the mask by one bit and create two new networks as a nodes left and right child. Left is the lower half, right the upper half.

## Roadmap

* network info (complete)
* basic subnetting (complete)
* summarization (complete)
* vlsm subnetting (in progress)
* full IPv6 support (tbd)

## Quick Start

### Compile and Install

```
go test ./...
go build
go install
```

