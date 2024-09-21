package ethernet

import (
	"encoding/binary"
	"errors"
)

const VIDMax = 0xFFF

// IEEE 802.1Q vlan structure
type VLAN struct {
	// Different PCP values can be used to prioritize different classes of traffic
	Priority uint8
	// This field was formerly designated Canonical Format Indicator (CFI)
	// with a value of 0 indicating a MAC address in canonical format.
	// It is always set to zero for Ethernet.
	DropEligible bool
	// Field specifying the VLAN to which the frame belongs.
	ID uint16
}

func (v *VLAN) Unmarshal(b []byte) error {
	if len(b) != 2 {
		return errors.New("size of vlan tag is 2 bytes")
	}

	num := binary.BigEndian.Uint16(b)

	v.Priority = uint8(num >> 13)
	v.DropEligible = num&0x1000 != 0
	v.ID = num & 0x0fff

	if v.ID > VIDMax {
		return errors.New("incorrect vlan identifier")
	}

	return nil
}

func (v *VLAN) Marshal() ([]byte, error) {
	b := make([]byte, 2)

	if v.ID > VIDMax {
		return nil, errors.New("incorrect vlan identifier")
	}

	num := uint16(v.Priority) << 13

	var drop uint16
	if v.DropEligible {
		drop = 1
	}

	num |= drop << 12
	num |= v.ID

	binary.BigEndian.PutUint16(b[:2], num)

	return b, nil
}
