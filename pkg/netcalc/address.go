package netcalc

import (
	"fmt"
	"strconv"
	"strings"
)

type IPv4Address struct {
	Address uint32
	Mask    uint32
}

type NetworkPrefix uint64

type IPv6Address struct {
	RoutingPrefix       uint64
	RoutingPrefixLength uint16
	SubnetId            uint16
	InterfaceIdentifier uint64
}

func (a *IPv6Address) NetworkPrefix() NetworkPrefix {
	var networkPrefix NetworkPrefix

	subnetLength := 64 - a.RoutingPrefixLength
	if subnetLength > 0 {
		networkPrefix = NetworkPrefix(a.RoutingPrefix << subnetLength)
		networkPrefix = networkPrefix | NetworkPrefix(a.SubnetId<<(16-subnetLength)>>(16-subnetLength))
	} else {
		networkPrefix = NetworkPrefix(a.RoutingPrefix)
	}

	return networkPrefix
}

func ParseAddress(addressString string) (uint32, error) {
	octets := strings.Split(addressString, ".")
	var address uint32 = 0
	if len(octets) == 4 {
		// correct number of octets
		for i, octet := range octets {
			// first octet is the highest 8 bits when shifted (24 bit shift)
			// fourth octet is the lower 8 bits when shifted (no shift)
			shift := (3 - i) * 8
			val, err := strconv.ParseUint(octet, 10, 32)

			if err != nil {
				return 0, err
			}

			val32 := uint32(val)

			if val32 > 255 {
				return 0, fmt.Errorf("utils:ParseAddress: parsing \"%s\": number must be 255 or less", octet)
			}

			// we have a good number
			address = address | val32<<shift

		}
	} else {
		// incorrect number of octets
		return 0, fmt.Errorf("utils:ParseAddress: parsing \"%s\": too many octets", addressString)
	}

	return address, nil
}

func ExportAddress(address uint32) string {

	firstOctet := address >> 24
	secondOctet := address << 8 >> 24
	thirdOctet := address << 16 >> 24
	fourthOctet := address << 24 >> 24

	ddAddress := strconv.FormatUint(uint64(firstOctet), 10) + "." + strconv.FormatUint(uint64(secondOctet), 10) + "." + strconv.FormatUint(uint64(thirdOctet), 10) + "." + strconv.FormatUint(uint64(fourthOctet), 10)

	return ddAddress
}

func GetNetworkAddress(address uint32, mask uint32) uint32 {
	return address & mask
}

func GetBroadcastAddress(address uint32, mask uint32) uint32 {
	return address | (^mask)
}

func IsValidMask(mask uint32) bool {
	for i := 1; i <= 32; i++ {
		calcMask, _ := GetMaskFromBits(i)
		if mask == calcMask {
			return true
		}
	}

	return false
}

func GetBitsInMask(mask uint32) int {
	bc := 0
	for mask != 0 {
		mask = mask << 1
		bc++
	}
	return bc
}

func GetMaskFromBits(bits int) (uint32, error) {
	if bits <= 32 {
		var mask uint32 = 0
		mask = ^mask
		bc := 32 - bits
		mask = mask << bc
		return mask, nil
	} else {
		return 0, fmt.Errorf("utils:GetMaskFromBits: bits must be 32 or less")
	}

}

func GetCommonBitMask(n1 uint32, n2 uint32) uint32 {
	commonBits := n1 ^ n2

	idx := 0

	for commonBits != 0 {
		commonBits = commonBits >> 1
		idx++
	}

	newMask, _ := GetMaskFromBits(32 - idx)

	return newMask
}
