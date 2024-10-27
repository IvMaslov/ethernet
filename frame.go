package ethernet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

// Most popular ether types
const (
	EtherTypeIPv4 = 0x0800
	EtherTypeARP  = 0x0806
	EtherTypeIPv6 = 0x86DD
	EtherTypeVLAN = 0x8100
	EtherTypeQinQ = 0x88a8
)

// Ethernet 802.3 frame
type Frame struct {
	// Destination MAC address
	DestHarwAddr net.HardwareAddr
	// Source MAC address
	SrcHarwAddr net.HardwareAddr
	// IEEE 802.1Q (VLAN) header
	// IEEE 802.1Q is the networking standard that supports virtual local area networking (VLANs)
	// on an IEEE 802.3 Ethernet network.
	VLAN *VLAN
	// IEEE 802.1ad (VLAN) header
	// QinQ frame is a frame that has two VLAN 802.1Q headers (i.e. it is double-tagged).
	ServiceVLAN *VLAN
	// EtherType is a two-octet field in an Ethernet frame.
	// It is used to indicate which protocol is encapsulated in the payload of the frame
	// https://en.wikipedia.org/wiki/EtherType
	EtherType uint16
	// Data of frame
	Payload []byte
}

func (f *Frame) Unmarshal(b []byte) error {
	if len(b) < 14 {
		return errors.New("frame is too short")
	}

	f.EtherType = binary.BigEndian.Uint16(b[12:14])

	payloadStart := 14

	switch f.EtherType {
	case EtherTypeVLAN:
		if len(b) < 18 {
			return errors.New("vlan frame is too short")
		}

		v := &VLAN{}

		err := v.Unmarshal(b[14:16])
		if err != nil {
			return fmt.Errorf("failed to unmarshal vlan header: %w", err)
		}

		f.VLAN = v
		f.EtherType = binary.BigEndian.Uint16(b[16:18])
		payloadStart = 18

	case EtherTypeQinQ:
		if len(b) < 22 {
			return errors.New("QinQ vlan frame is too short")
		}

		serviceTag := &VLAN{}
		customerTag := &VLAN{}

		err := serviceTag.Unmarshal(b[14:16])
		if err != nil {
			return fmt.Errorf("failed to unmarshal QinQ vlan header: %w", err)
		}

		err = customerTag.Unmarshal(b[18:20])
		if err != nil {
			return fmt.Errorf("failed to unmarshal vlan header: %w", err)
		}

		f.VLAN = customerTag
		f.ServiceVLAN = serviceTag
		f.EtherType = binary.BigEndian.Uint16(b[20:22])
		payloadStart = 22
	}

	buf := make([]byte, 12+len(b[payloadStart:]))

	copy(buf[:6], b[:6])
	copy(buf[6:12], b[6:12])
	copy(buf[12:], b[payloadStart:])

	f.DestHarwAddr = buf[:6]
	f.SrcHarwAddr = buf[6:12]
	f.Payload = buf[12:]

	return nil
}

func (f *Frame) Marshal() ([]byte, error) {
	buf := make([]byte, f.length())

	copy(buf[:6], f.DestHarwAddr)
	copy(buf[6:12], f.SrcHarwAddr)

	payloadStart := 14

	if f.ServiceVLAN != nil {
		binary.BigEndian.PutUint16(buf[12:14], EtherTypeQinQ)

		v, err := f.ServiceVLAN.Marshal()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal QinQ vlan header: %w", err)
		}

		copy(buf[14:16], v)

		payloadStart = 18
	}

	if f.VLAN != nil {
		binary.BigEndian.PutUint16(buf[payloadStart-2:payloadStart], EtherTypeVLAN)

		v, err := f.VLAN.Marshal()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal vlan header: %w", err)
		}

		copy(buf[payloadStart:payloadStart+2], v)

		payloadStart = payloadStart + 4
	}

	binary.BigEndian.PutUint16(buf[payloadStart-2:payloadStart], f.EtherType)

	copy(buf[payloadStart:], f.Payload)

	return buf, nil
}

func (f *Frame) length() int {
	l := 0

	l += 12 // two MAC addresses
	l += 2  // EtherType

	if f.ServiceVLAN != nil {
		l += 4
	}

	if f.VLAN != nil {
		l += 4
	}

	return l + len(f.Payload)
}
