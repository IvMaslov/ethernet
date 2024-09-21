package ethernet

import (
	"testing"
)

func Test_Frame_Unmarshal(t *testing.T) {
	tests := []struct {
		input    []byte
		expected Frame
	}{
		{
			input: []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x08, 0x00, 0x50, 0x50, 0x50},
			expected: Frame{
				DestHarwAddr: []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
				SrcHarwAddr:  []byte{0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
				EtherType:    EtherTypeIPv4,
				Payload:      []byte{0x50, 0x50, 0x50},
			},
		},
		{
			input: []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x81, 0x00, 0x70, 0x0c, 0x08, 0x00, 0x50, 0x50, 0x50},
			expected: Frame{
				DestHarwAddr: []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
				SrcHarwAddr:  []byte{0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
				EtherType:    EtherTypeIPv4,
				VLAN: &VLAN{
					Priority:     3,
					ID:           12,
					DropEligible: true,
				},
				Payload: []byte{0x50, 0x50, 0x50},
			},
		},
		{
			input: []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x88, 0xa8, 0x70, 0x0c, 0x81, 0x00, 0x70, 0x0c, 0x08, 0x00, 0x50, 0x50, 0x50},
			expected: Frame{
				DestHarwAddr: []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
				SrcHarwAddr:  []byte{0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
				EtherType:    EtherTypeIPv4,
				VLAN: &VLAN{
					Priority:     3,
					ID:           12,
					DropEligible: true,
				},
				ServiceVLAN: &VLAN{
					Priority:     3,
					ID:           12,
					DropEligible: true,
				},
				Payload: []byte{0x50, 0x50, 0x50},
			},
		},
	}

	for _, test := range tests {
		f := &Frame{}

		err := f.Unmarshal(test.input)
		if err != nil {
			t.Errorf("failed to unmarshal: %v", err)
		}

		for i, b := range test.expected.DestHarwAddr {
			if b != f.DestHarwAddr[i] {
				t.Error("wrong destination address")
			}
		}

		for i, b := range test.expected.SrcHarwAddr {
			if b != f.SrcHarwAddr[i] {
				t.Error("wrong source address")
			}
		}

		if test.expected.EtherType != f.EtherType {
			t.Error("wrong ethertype")
		}

		for i, b := range test.expected.Payload {
			if b != f.Payload[i] {
				t.Error("wrong payload")
			}
		}

		if test.expected.VLAN != nil &&
			test.expected.VLAN.DropEligible != f.VLAN.DropEligible &&
			test.expected.VLAN.Priority != f.VLAN.Priority &&
			test.expected.VLAN.ID != f.VLAN.ID {
			t.Error("wrong VLAN header")
		}

		if test.expected.ServiceVLAN != nil &&
			test.expected.ServiceVLAN.DropEligible != f.ServiceVLAN.DropEligible &&
			test.expected.ServiceVLAN.Priority != f.ServiceVLAN.Priority &&
			test.expected.ServiceVLAN.ID != f.ServiceVLAN.ID {
			t.Error("wrong QinQ VLAN header")
		}
	}
}

func Test_Frame_Marshal(t *testing.T) {
	tests := []struct {
		f        Frame
		expected []byte
	}{
		{
			f: Frame{
				DestHarwAddr: []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
				SrcHarwAddr:  []byte{0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
				EtherType:    EtherTypeIPv4,
				Payload:      []byte{0x50, 0x50, 0x50},
			},
			expected: []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x08, 0x00, 0x50, 0x50, 0x50},
		},
		{
			f: Frame{
				DestHarwAddr: []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
				SrcHarwAddr:  []byte{0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
				EtherType:    EtherTypeIPv4,
				VLAN: &VLAN{
					Priority:     3,
					ID:           12,
					DropEligible: true,
				},
				Payload: []byte{0x50, 0x50, 0x50},
			},
			expected: []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x81, 0x00, 0x70, 0x0c, 0x08, 0x00, 0x50, 0x50, 0x50},
		},
		{
			f: Frame{
				DestHarwAddr: []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
				SrcHarwAddr:  []byte{0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
				EtherType:    EtherTypeIPv4,
				VLAN: &VLAN{
					Priority:     3,
					ID:           12,
					DropEligible: true,
				},
				ServiceVLAN: &VLAN{
					Priority:     3,
					ID:           12,
					DropEligible: true,
				},
				Payload: []byte{0x50, 0x50, 0x50},
			},
			expected: []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x88, 0xa8, 0x70, 0x0c, 0x81, 0x00, 0x70, 0x0c, 0x08, 0x00, 0x50, 0x50, 0x50},
		},
	}

	for _, test := range tests {
		buf, err := test.f.Marshal()
		if err != nil {
			t.Error("failed to marshal")
		}

		for i, b := range test.expected {
			if b != buf[i] {
				t.Error("wrong frame data")
			}
		}
	}
}