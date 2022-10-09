package utils

import (
	"fmt"
	"strconv"
	"strings"
)

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
