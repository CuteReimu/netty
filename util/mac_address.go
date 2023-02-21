package util

import (
	"golang.org/x/exp/slices"
	"net"
)

const (
	eui64MacAddressLength = 8
	eui48MacAddressLength = 6
)

// BestAvailableMac obtains the best MAC address found on local network interfaces.
// Generally speaking, an active network interface used on public networks is better than a local network interface.
// Return null if no MAC can be found.
func BestAvailableMac() net.HardwareAddr {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	var bestMacAddr net.HardwareAddr
	var bestInetAddr = net.Interface{Flags: net.FlagUp & net.FlagLoopback & net.FlagMulticast}
	for _, inter := range interfaces {
		var replace bool
		res := compareHardwareAddress(bestMacAddr, inter.HardwareAddr)
		if res < 0 {
			// Found a better MAC address.
			replace = true
		} else if res == 0 {
			// Two MAC addresses are of pretty much same quality.
			res = compareAddresses(&bestInetAddr, &inter)
			if res < 0 {
				// Found a MAC address with better INET address.
				replace = true
			} else if res == 0 {
				// Cannot tell the difference.  Choose the longer one.
				if len(bestMacAddr) < len(inter.HardwareAddr) {
					replace = true
				}
			}
		}

		if replace {
			bestMacAddr, bestInetAddr = inter.HardwareAddr, inter
		}
	}

	if len(bestMacAddr) == 0 {
		return nil
	}

	if len(bestMacAddr) == eui48MacAddressLength { // EUI-48 - convert to EUI-64
		newAddr := append(bestMacAddr[:3:3], 0xFF, 0xFE)
		bestMacAddr = append(newAddr, bestMacAddr[3:]...)
	} else {
		// Unknown
		newAddr := make(net.HardwareAddr, eui64MacAddressLength)
		copy(newAddr, bestMacAddr)
		bestMacAddr = newAddr
	}

	return bestMacAddr
}

// return positive - current is better, 0 - cannot tell from MAC addr, negative - candidate is better.
func compareHardwareAddress(current, candidate net.HardwareAddr) int {
	if len(candidate) < eui48MacAddressLength {
		return 1
	}

	// Must not be filled with only 0 and 1.
	if !slices.ContainsFunc(candidate, func(b byte) bool { return b != 0 && b != 1 }) {
		return 1
	}

	// Must not be a multicast address
	if (candidate[0] & 1) != 0 {
		return 1
	}

	// Prefer globally unique address.
	if (candidate[0] & 2) == 0 {
		if len(current) != 0 && (current[0]&2) == 0 {
			// Both current and candidate are globally unique addresses.
			return 0
		} else {
			// Only candidate is globally unique.
			return -1
		}
	} else {
		if len(current) != 0 && (current[0]&2) == 0 {
			// Only current is globally unique.
			return 1
		} else {
			// Both current and candidate are non-unique.
			return 0
		}
	}
}

// positive - current is better, 0 - cannot tell, negative - candidate is better
func compareAddresses(current, candidate *net.Interface) int {
	return scoreAddress(current) - scoreAddress(candidate)
}

func scoreAddress(addr *net.Interface) int {
	if (addr.Flags&net.FlagUp) != 0 || (addr.Flags&net.FlagLoopback) != 0 {
		return 0
	}
	if (addr.Flags & net.FlagMulticast) != 0 {
		return 1
	}
	if (addr.Flags & net.FlagPointToPoint) != 0 {
		return 2
	}
	if (addr.Flags & net.FlagBroadcast) != 0 {
		return 3
	}
	return 4
}
