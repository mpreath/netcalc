# netcalc

## Overview
netcalc is an advanced CLI-based network subnetting tool. It provides the
following features:

* provided an ip address and mask, it will calculate the network address, broadcast address, bits in mask, and all hosts of the network
* provided an ip address, mask, and network count, it will calculate the network addresses, and masks in order to subnet the network into smaller networks based on the number of networks desired
* provided an ip address, mask, and host count, it will calculate the network addresses, and mask in ord to subnet the network into smaller networks based on the number of hosts desired 
* provided an ip address, mask, and a list of host count values, it will calculate a list of VLSM subnets that match the list of host count requirements. Instead of making the subnets all the same size, it only creates enough subnets to match the list, and they may all be different sizes
* provided a list of subnets, it will try to summarize the subnets into one (or more) supernets if possible.
* 

*NOTE:* netcalc uses bit shifting and mask operators to work with ip addresses in their actual 32 bit integer form. It converts from dotted decimal notation to unsigned integer values, and then manipulates the integer values.

*NOTE:* netcalc uses a binary tree structure to subnet networks by extending the mask by one bit and create two new networks as a nodes left and right child. Left is the lower half, right the upper half.

## Roadmap

1. network info (in progress)
2. basic subnetting (tbd)
3. vlsm subnetting (tbd)
4. supernetting (tbd)
5. network info ipv6 (tbd)

## Quick Start

### Compile and Install

```
go build
go install
```

### Usage

```
netcalc [OPTION?] [<ip_address/cidr> | <ip_address> <mask>] - calculate network information

Help Options:
  -h, --help               Show help options

Application Options:
  -v, --verbose            Use verbose output
  -V, --version            Display version information
  -o, --hosts=HOSTS        Specify the host count to use to subnet the network
  -n, --nets=NETS          Specify the net count to use to subnet the network
  -l, --vlsm=VLSM_LIST     Comma seperated list of host counts for VLSM network
  -s, --summary            Summarize subnets into one or more supernets
 ```

