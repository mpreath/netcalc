package utils

import (
	"errors"
)

func GetNetworkAddress(address uint32, mask uint32) uint32 {
	return address & mask
}

func GetBroadcastAddress(address uint32, mask uint32) uint32 {
	return address | (^mask)
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
