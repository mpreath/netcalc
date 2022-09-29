![Netcalc Logo](https://github.com/mpreath/public_files/blob/25a6d05318b3ddfbf7b3d3835201aa1a0cfa48ab/netcalc.png)

## Overview
netcalc is an advanced CLI-based network subnetting tool. It provides the
following features:

* provided an ip address and mask, it will calculate the network address, broadcast address, bits in mask, and all hosts of the network
* provided an ip address, mask, and network count, it will calculate the network addresses, and masks in order to subnet the network into smaller networks based on the number of networks desired
* provided an ip address, mask, and host count, it will calculate the network addresses, and mask in ord to subnet the network into smaller networks based on the number of hosts desired 
* provided an ip address, mask, and a list of host count values, it will calculate a list of VLSM subnets that match the list of host count requirements. Instead of making the subnets all the same size, it only creates enough subnets to match the list, and they may all be different sizes
* provided a list of subnets, it will try to summarize the subnets into one (or more) supernets if possible.

## Todo

* Refactor data structure to use less memory (in progress)
* Implement IPv6 support (on roadmap)

## Install

```
go test ./...
go build
go install
```
## Usage
```
Netcalc is a IPv4/IPv6 network calculator

Usage:
  netcalc [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  info        Displays information about a network
  subnet      Given a network break it into smaller networks
  summarize   summarizes the networks provided to stdin
  version     Print the version number of netcalc
  vlsm        Given a network and comma-separated list of subnet lengths break it into smaller networks

Flags:
  -h, --help      help for netcalc
  -j, --json      Turns on JSON output for commands
  -v, --verbose   Turns on verbose output for commands

Use "netcalc [command] --help" for more information about a command.
```
### `info` Command
```
Usage:
  netcalc info <ip_address> <subnet_mask> [flags]

Flags:
  -h, --help   help for info

Global Flags:
  -j, --json      Turns on JSON output for commands
  -v, --verbose   Turns on verbose output for commands
```

```
> netcalc info 192.168.0.0 255.255.255.248 -v
Network:        192.168.0.0
Mask:           255.255.255.248 (/29)
Bcast:          192.168.0.7

192.168.0.1     255.255.255.248
192.168.0.2     255.255.255.248
192.168.0.3     255.255.255.248
192.168.0.4     255.255.255.248
192.168.0.5     255.255.255.248
192.168.0.6     255.255.255.248
```
### `subnet` Command

```
Usage:
  netcalc subnet [--hosts <hosts> | --networks <networks>] <ip_address> <subnet_mask> [flags]

Flags:
  -h, --help           help for subnet
      --hosts int      Specifies the number of hosts to include each subnet.
      --networks int   Specifies the number of subnets to create.

Global Flags:
  -j, --json      Turns on JSON output for commands
  -v, --verbose   Turns on verbose output for commands
```

```
> netcalc subnet --hosts 2 192.168.1.0 255.255.255.224
192.168.1.0     255.255.255.252
192.168.1.4     255.255.255.252
192.168.1.8     255.255.255.252
192.168.1.12    255.255.255.252
192.168.1.16    255.255.255.252
192.168.1.20    255.255.255.252
192.168.1.24    255.255.255.252
192.168.1.28    255.255.255.252
```

```
> netcalc subnet -v --hosts 2 192.168.1.0 255.255.255.224
* = assigned network
+ = useable network
[n] = # of useable hosts

__192.168.1.0/27
 |__192.168.1.0/28
 | |__192.168.1.0/29
 | | |__192.168.1.0/30[2]+
 | | |__192.168.1.4/30[2]+
 | |__192.168.1.8/29
 | | |__192.168.1.8/30[2]+
 | | |__192.168.1.12/30[2]+
 |__192.168.1.16/28
 | |__192.168.1.16/29
 | | |__192.168.1.16/30[2]+
 | | |__192.168.1.20/30[2]+
 | |__192.168.1.24/29
 | | |__192.168.1.24/30[2]+
 | | |__192.168.1.28/30[2]+
```

```
> netcalc subnet -j --hosts 2 192.168.1.0 255.255.255.248
{
  "network": {
    "address": "192.168.1.0",
    "mask": "255.255.255.248",
    "broadcast": "192.168.1.7"
  },
  "subnets": [
    {
      "network": {
        "address": "192.168.1.0",
        "mask": "255.255.255.252",
        "broadcast": "192.168.1.3"
      }
    },
    {
      "network": {
        "address": "192.168.1.4",
        "mask": "255.255.255.252",
        "broadcast": "192.168.1.7"
      }
    }
  ]
}
```

### `vlsm` Command

```
Usage:
  netcalc vlsm <host_counts_list> <ip_address> <subnet_mask> [flags]

Flags:
  -h, --help   help for vlsm

Global Flags:
  -j, --json      Turns on JSON output for commands
  -v, --verbose   Turns on verbose output for commands
```

The `subnet` command above splits the network provided into a smaller of number of networks based on either host count or number of networks desired, all networks created have the same length subnet mask.  The `vlsm` command also splits the network provided into a smaller number of networks but rather lets you provide a comma-separated list of host counts.  

For example, for this example we start with `192.168.1.0 255.255.255.128` as our base network and have a need to subnet this further to support: 
1. Four (4) /30 point-to-point links (links between routers)
2. Two (2) 20 host count networks (two small office networks)

You can accomplish this with `netcalc` by specifying the desired host count values as arguments to the command. 

```
> netcalc vlsm 2,2,2,2,20,20 192.168.1.0 255.255.255.128
192.168.1.0     255.255.255.252
192.168.1.4     255.255.255.252
192.168.1.8     255.255.255.252
192.168.1.12    255.255.255.252
192.168.1.16    255.255.255.240
192.168.1.32    255.255.255.224
192.168.1.64    255.255.255.224
192.168.1.96    255.255.255.224
```

```
netcalc vlsm -v 2,2,2,2,20,20 192.168.1.0 255.255.255.128
* = assigned network
+ = useable network
[n] = # of useable hosts

__192.168.1.0/25
 |__192.168.1.0/26
 | |__192.168.1.0/27
 | | |__192.168.1.0/28
 | | | |__192.168.1.0/29
 | | | | |__192.168.1.0/30[2]*
 | | | | |__192.168.1.4/30[2]*
 | | | |__192.168.1.8/29
 | | | | |__192.168.1.8/30[2]*
 | | | | |__192.168.1.12/30[2]*
 | | |__192.168.1.16/28[14]+
 | |__192.168.1.32/27[30]*
 |__192.168.1.64/26
 | |__192.168.1.64/27[30]*
 | |__192.168.1.96/27[30]+
```

### `summarize` Command

```
summarizes the networks provided to stdin

Usage:
  netcalc summarize [flags]

Flags:
  -h, --help   help for summarize

Global Flags:
  -j, --json      Turns on JSON output for commands
  -v, --verbose   Turns on verbose output for command
```

The summarize command takes a list of networks from `<stdin>` and summarizes them at the bit level into a single summary network.  This is useful for determining a common network route to configure in a routing environment to simplify routing advertisements. 

In this example there is a text file (`test/vlsm_networks.txt`) containing the following networks:
```
192.168.1.0     255.255.255.252
192.168.1.4     255.255.255.252
192.168.1.8     255.255.255.252
192.168.1.12    255.255.255.252
192.168.1.16    255.255.255.240
192.168.1.32    255.255.255.224
192.168.1.64    255.255.255.224
192.168.1.96    255.255.255.224
```

```
> netcalc summarize < test/vlsm_networks.txt                                     
192.168.1.0     255.255.255.128
```

The summarize command takes those networks and summarizes at the bit level to create the summary network of `192.168.1.0 255.255.255.128`. 

