package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Ddtoi(address_string string) (uint32, error) {
	octets := strings.Split(address_string, ".")
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
				return 0, fmt.Errorf("ipv4:Ddtoi: parsing \"%s\": number must be 255 or less", octet)
			}

			// we have a good number
			address = address | val32<<shift

		}
	} else {
		// incorrect number of octets
		return 0, fmt.Errorf("ipv4:Ddtoi: parsing \"%s\": too many octets", address_string)
	}

	return address, nil
}

func Itodd(address uint32) string {

	first_octet := address >> 24
	second_octet := address << 8 >> 24
	third_octet := address << 16 >> 24
	fourth_octet := address << 24 >> 24

	dd_address := strconv.FormatUint(uint64(first_octet), 10) + "." + strconv.FormatUint(uint64(second_octet), 10) + "." + strconv.FormatUint(uint64(third_octet), 10) + "." + strconv.FormatUint(uint64(fourth_octet), 10)

	return dd_address
}

func IsValidMask(mask uint32) bool {
	for i := 1; i <= 32; i++ {
		calc_mask, _ := GetMaskFromBits(i)
		if mask == calc_mask {
			return true
		}
	}

	return false
}

func GetNetworkAddress(address uint32, mask uint32) uint32 {
	return address & mask
}

func GetBroadcastAddress(address uint32, mask uint32) uint32 {
	return address | (^mask)
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
		return 0, errors.New("ipv4:GetMaskFromBits: bits must be 32 or less")
	}

}
