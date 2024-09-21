package ethernet

import (
	"testing"
)

func Test_VLAN_Unmarshal(t *testing.T) {
	tests := []struct {
		input    []byte
		expected VLAN
	}{
		{
			input: []byte{0x70, 0x0c},
			expected: VLAN{
				Priority:     3,
				ID:           12,
				DropEligible: true,
			},
		},
		{
			input: []byte{0xa0, 0xcc},
			expected: VLAN{
				Priority: 5,
				ID:       204,
			},
		},
	}

	for _, test := range tests {
		var v VLAN

		err := v.Unmarshal(test.input)

		if err != nil {
			t.Errorf("unexpected err %v", err)
		}

		if v.DropEligible != test.expected.DropEligible || v.ID != test.expected.ID || v.Priority != test.expected.Priority {
			t.Error("error while unmarshalling")
		}
	}
}

func Test_VLAN_Marshal(t *testing.T) {
	tests := []struct {
		v        VLAN
		expected []byte
	}{
		{
			v: VLAN{
				Priority:     3,
				ID:           12,
				DropEligible: true,
			},
			expected: []byte{0x70, 0x0c},
		},
		{
			v: VLAN{
				Priority: 5,
				ID:       204,
			},
			expected: []byte{0xa0, 0xcc},
		},
	}

	for _, test := range tests {
		buf, err := test.v.Marshal()

		if err != nil {
			t.Errorf("unexpected err %v", err)
		}

		for i, b := range test.expected {
			if b != buf[i] {
				t.Error("error while marshalling")
			}
		}
	}
}
