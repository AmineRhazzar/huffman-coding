package huffman

import (
	"fmt"
	"io"
	"testing"
)

func TestReadBit_ShouldFail(t *testing.T) {
	type test_case struct {
		description string
		data        []byte
		cursor      uint8
		idx         int
	}

	data := []byte{0b0000_0000, 0b0100_0101, 0b0000_0101}
	test_cases := []test_case{
		{
			description: "index at last byte",
			data:        data,
			cursor:      1,
			idx:         len(data) - 1,
		},
		{
			description: "index out of bounds",
			data:        data, cursor: 1,
			idx: 5,
		},
		{
			description: "cursor past read size specified in last byte",
			data:        data, cursor: 5,
			idx: 1,
		},
	}

	for scenarioIdx, scenario := range test_cases {
		t.Run(fmt.Sprintf("Test %d", scenarioIdx), func(t *testing.T) {
			r := GetReader(scenario.data)
			r.cursor = scenario.cursor
			r.idx = scenario.idx

			_, err := r.ReadBit()

			if err != io.EOF {
				t.Fatalf(`Test %d Failed.
				Got:err %v,
				Wanted: err %v`, scenarioIdx, err, io.EOF)
			}
		})
	}
}

func TestReadByte_ShouldFail(t *testing.T) {
	type test_case struct {
		description string
		cursor      uint8
		idx         int
		data        []byte
	}

	test_cases := []test_case{
		{
			description: "index out of bounds",
			data:        []byte{0b0000_0000, 0b0100_0101, 0b0000_0101},
			idx:         4,
		}, {
			description: "index at last byte",
			data:        []byte{0b0000_0000, 0b0100_0101, 0b0000_0101},
			idx:         2,
		}, {
			description: "index at second last byte, second_last_byte_read_size != 0 ",
			data:        []byte{0b0000_0000, 0b0100_0101, 0b0000_0101},
			idx:         1,
		},
		{
			description: "index at second last byte, second_last_byte_read_size = 0, cursor != zero",
			data:        []byte{0b0000_0000, 0b0100_0101, 0b0000_0101, 0b0000_0000},
			cursor:      1,
			idx:         2,
		},
		{
			description: "read byte reaching second last byte, reading past second_last_byte_read_size",
			// in this case the last byte is 2 so we can only read up to here
			//                                                 v
			data: []byte{0b0011_0010, 0b0100_0101, 0b0000_0101, 0b0000_0010},
			// should fail because we start from here ^..........^ cursor goes more than read size limit
			idx:    1,
			cursor: 4,
		},
	}

	for scenarioIdx, scenario := range test_cases {
		t.Run(scenario.description, func(t *testing.T) {
			r := GetReader(scenario.data)
			r.cursor = scenario.cursor
			r.idx = scenario.idx

			b, err := r.ReadByte()

			if err == nil {
				t.Fatalf(`Test %d Failed.
				Got: %08b, %v,
				Wanted:0, <nil>`, scenarioIdx, b, err)
			}
		})
	}
}

func TestReadBit(t *testing.T) {
	type test_case struct {
		description  string
		cursor       uint8
		idx          int
		data         []byte
		expected_bit uint8
	}
	test_cases := []test_case{
		{
			description:  "read bit",
			data:         []byte{0b0011_0010, 0b0100_0101, 0b0000_0101, 0b0000_0000},
			idx:          1,
			cursor:       2,
			expected_bit: 0,
		},
		{
			description:  "read bit on second last byte, second_last_byte_read_size=0",
			data:         []byte{0b0011_0010, 0b0100_0101, 0b0000_0101, 0b0000_0000},
			idx:          2,
			cursor:       7,
			expected_bit: 1,
		},
		{
			description:  "read bit on second last byte, second_last_byte_read_size != 0",
			data:         []byte{0b0011_0010, 0b0100_0101, 0b0000_0101, 0b0000_0101},
			idx:          2,
			cursor:       4,
			expected_bit: 0,
		},
	}

	for scenarioIdx, scenario := range test_cases {
		t.Run(scenario.description, func(t *testing.T) {
			r := GetReader(scenario.data)
			r.cursor = scenario.cursor
			r.idx = scenario.idx

			bit, err := r.ReadBit()

			if err != nil || scenario.expected_bit != bit {
				t.Fatalf(`Test %d Failed.
				Got: %d, %v,
				Wanted: %d, <nil>`, scenarioIdx, bit, err, scenario.expected_bit)
			}
		})
	}
}

func TestReadByte(t *testing.T) {
	type test_case struct {
		description   string
		cursor        uint8
		idx           int
		data          []byte
		expected_byte uint8
	}
	test_cases := []test_case{
		{
			description:   "read byte, cursor = 0",
			data:          []byte{0b0011_0010, 0b0100_0101, 0b0000_0101, 0b0000_0000},
			idx:           1,
			cursor:        0,
			expected_byte: 0b0100_0101,
		},
		{
			description:   "read byte, cursor != 0",
			data:          []byte{0b0011_0010, 0b0100_0101, 0b0000_0101, 0b0101_1100, 0b0000_0000},
			idx:           1,
			cursor:        2,
			expected_byte: 0b0001_0100,
		},
		{
			description:   "read byte reaching second last byte, second_last_byte_read_size = 0",
			data:          []byte{0b0011_0010, 0b0100_0101, 0b0000_0101, 0b0000_0000},
			idx:           1,
			cursor:        6,
			expected_byte: 0b0100_0001,
		},
		{
			description:   "read byte reaching second last byte, second_last_byte_read_size != 0",
			data:          []byte{0b0011_0010, 0b0100_0101, 0b0000_0101, 0b0000_0100},
			idx:           1,
			cursor:        4,
			expected_byte: 0b01010000,
		},
	}

	for scenarioIdx, scenario := range test_cases {
		t.Run(scenario.description, func(t *testing.T) {
			r := GetReader(scenario.data)
			r.cursor = scenario.cursor
			r.idx = scenario.idx

			b, err := r.ReadByte()

			if err != nil || scenario.expected_byte != b || r.cursor != scenario.cursor || r.idx != scenario.idx+1 {
				t.Fatalf(`Test %d Failed.
				Got: %08b, %v, (idx=%d, cursor=%d)
				Wanted: %08b, <nil>, (idx=%d, cursor=%d)`, scenarioIdx, b, err, r.idx, r.cursor, scenario.expected_byte, scenario.idx+1, scenario.cursor)
			}
		})
	}
}

func TestReadTree(t *testing.T) {
	type test_case struct {
		description string
		cursor      uint8
		idx         int
		data        []byte
	}

	test_cases := []test_case{
		{
			description: "read tree, index=0,cursor=0",
			data:        []byte{0b00101000, 0b01110100, 0b00010101, 0b00010101, 0b01000100, 0b10100001, 0b00000000, 0b00000001},
			idx:         0,
			cursor:      0,
		},
		{
			description: "read tree, index!=0, cursor != 0,",
			data:        []byte{0b00110001, 0b01001010, 0b00011101, 0b00000101, 0b01000101, 0b01010001, 0b00101000, 0b01000000, 0b00000011},
			//                                  ^10                                                                      ^59
			idx:    1,
			cursor: 2,
		},
		{
			description: "read tree, buffer bigger than tree",
			data:        []byte{0b00110001, 0b01001010, 0b00011101, 0b00000101, 0b01000101, 0b01010001, 0b00101000, 0b01010010, 0b00100011, 0b00101010, 0b0000_0101},
			idx:         1,
			cursor:      2,
		},
	}

	for scenarioIdx, scenario := range test_cases {
		t.Run(scenario.description, func(t *testing.T) {
			r := GetReader(scenario.data)
			r.cursor = scenario.cursor
			r.idx = scenario.idx

			start_index := r.idx*8 + int(r.cursor)
			tree, err := r.ReadTree(start_index, MOCK_TREE_ENCODED_SIZE)
			current_index := r.idx*8 + int(r.cursor)

			equal_trees := areTreesEqual(tree, MOCK_TREE)
			if err != nil ||
				current_index-start_index != MOCK_TREE_ENCODED_SIZE ||
				!equal_trees {
				t.Fatalf(`Test %d Failed.
				Got: equal trees: %v, err: %v, size read: %v
				Wanted: equal trees: true, err: <nil>, size read: %d`, scenarioIdx, equal_trees, err, current_index-start_index, MOCK_TREE_ENCODED_SIZE)
			}
		})
	}

}
