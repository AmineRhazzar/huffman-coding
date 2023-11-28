package huffman

import (
	"testing"
)

/*
  - 20
    -----------------
    12				8

-----------     ---------
C,6		A,6		E,5		3

	-----------
	D,2		B,1
*/
var MOCK_TREE *Node = &Node{
	weight: 20,
	Left: &Node{
		weight: 12,
		Left:   &Node{ch: 'C', weight: 6},
		Right:  &Node{ch: 'A', weight: 6},
	},
	Right: &Node{
		weight: 8,
		Left:   &Node{ch: 'E', weight: 5},
		Right: &Node{weight: 3,
			Left:  &Node{ch: 'D', weight: 2},
			Right: &Node{ch: 'B', weight: 1},
		},
	},
}

const MOCK_TREE_ENCODED_SIZE int = 49

func TestWriteBit(t *testing.T) {
	type test_case struct {
		description     string
		bit             uint8
		writer          Writer
		expected_writer Writer
	}

	test_cases := []test_case{
		{
			description: "write bit 1 on empty buffer, cursor at 0",
			bit:         1,
			writer: Writer{
				buffer:    []byte{},
				curr_byte: 0b0000_0000,
				cursor:    0,
			},
			expected_writer: Writer{
				buffer:    []byte{},
				curr_byte: 0b1000_0000,
				cursor:    1,
			},
		},
		{
			description: "write bit 0 on empty buffer, cursor at 0",
			bit:         0,
			writer: Writer{
				buffer:    []byte{},
				curr_byte: 0b0000_0000,
				cursor:    0,
			},
			expected_writer: Writer{
				buffer:    []byte{},
				curr_byte: 0b0000_0000,
				cursor:    1,
			},
		},
		{
			description: "write bit 1 on empty buffer, cursor at 3",
			bit:         1,
			writer: Writer{
				buffer:    []byte{},
				curr_byte: 0b1000_0000,
				cursor:    3,
			},
			expected_writer: Writer{
				buffer:    []byte{},
				curr_byte: 0b1001_0000,
				cursor:    4,
			},
		},
		{
			description: "write bit 0 on empty buffer, cursor at 7",
			bit:         0,
			writer: Writer{
				buffer:    []byte{},
				curr_byte: 0b1000_0000,
				cursor:    7,
			},
			expected_writer: Writer{
				buffer:    []byte{0b1000_0000},
				curr_byte: 0b0000_0000,
				cursor:    0,
			},
		},
		{
			description: "write bit 1 on non-empty buffer, cursor at 7",
			bit:         1,
			writer: Writer{
				buffer:    []byte{0b1000_0000},
				curr_byte: 0b1010_0010,
				cursor:    7,
			},
			expected_writer: Writer{
				buffer:    []byte{0b1000_0000, 0b1010_0011},
				curr_byte: 0b0000_0000,
				cursor:    0,
			},
		},
	}

	for scenarioIdx, scenario := range test_cases {
		t.Run(scenario.description, func(t *testing.T) {
			scenario.writer.WriteBit(scenario.bit)
			if !areWritersEqual(scenario.writer, scenario.expected_writer) {
				t.Fatalf(`Test %d Failed.
				writer: %+v,
				expected: %+v`, scenarioIdx, scenario.writer, scenario.expected_writer)
			}
		})
	}
}

func TestWriteMultipleBits(t *testing.T) {
	type test_case struct {
		description     string
		bits            []uint8
		writer          Writer
		expected_writer Writer
	}

	test_cases := []test_case{
		{
			description: "write bits 101 on empty buffer, cursor at 0",
			bits:        []uint8{1, 0, 1},
			writer: Writer{
				buffer:    []byte{},
				curr_byte: 0b0000_0000,
				cursor:    0,
			},
			expected_writer: Writer{
				buffer:    []byte{},
				curr_byte: 0b1010_0000,
				cursor:    3,
			},
		},
		{
			description: "write bits 0101 on empty buffer, cursor at 4",
			bits:        []uint8{0, 1, 0, 1},
			writer: Writer{
				buffer:    []byte{},
				curr_byte: 0b0000_0000,
				cursor:    4,
			},
			expected_writer: Writer{
				buffer:    []byte{0b0000_0101},
				curr_byte: 0b0000_0000,
				cursor:    0,
			},
		},
		{
			description: "write bits 110000011 on empty buffer, cursor at 3",
			bits:        []uint8{1, 1, 0, 0, 0, 0, 0, 1, 1},
			writer: Writer{
				buffer:    []byte{},
				curr_byte: 0b1000_0000,
				cursor:    3,
			},
			expected_writer: Writer{
				buffer:    []byte{0b1001_1000},
				curr_byte: 0b0011_0000,
				cursor:    4,
			},
		},
	}

	for scenarioIdx, scenario := range test_cases {
		t.Run(scenario.description, func(t *testing.T) {
			scenario.writer.WriteMultipleBits(scenario.bits...)
			if !areWritersEqual(scenario.writer, scenario.expected_writer) {
				t.Fatalf(`Test %d Failed.
				writer: %+v,
				expected: %+v`, scenarioIdx, scenario.writer, scenario.expected_writer)
			}
		})
	}
}

func TestWriteByte(t *testing.T) {
	type test_case struct {
		description     string
		b               byte
		writer          Writer
		expected_writer Writer
	}

	test_cases := []test_case{
		{
			description: "write byte 01010001 on empty buffer, cursor at 0",
			b:           0b01010001,
			writer: Writer{
				buffer:    []byte{},
				curr_byte: 0b0000_0000,
				cursor:    0,
			},
			expected_writer: Writer{
				buffer:    []byte{0b01010001},
				curr_byte: 0b0000_0000,
				cursor:    0,
			},
		},
		{
			description: "write bits 01010001 on empty buffer, cursor at 4",
			b:           0b01010001,
			writer: Writer{
				buffer:    []byte{},
				curr_byte: 0b0100_0000,
				cursor:    4,
			},
			expected_writer: Writer{
				buffer:    []byte{0b0100_0101},
				curr_byte: 0b0001_0000,
				cursor:    4,
			},
		},
	}

	for scenarioIdx, scenario := range test_cases {
		t.Run(scenario.description, func(t *testing.T) {
			scenario.writer.WriteByte(scenario.b)
			if !areWritersEqual(scenario.writer, scenario.expected_writer) {
				t.Fatalf(`Test %d Failed.
				writer: %+v,
				expected: %+v`, scenarioIdx, scenario.writer, scenario.expected_writer)
			}
		})
	}
}

func TestWriteTree(t *testing.T) {
	w := Writer{}

	tree_encoded_size := w.WriteTree(MOCK_TREE)
	expected_writer := Writer{
		buffer:    []byte{0b00101000, 0b01110100, 0b00010101, 0b00010101, 0b01000100, 0b10100001},
		cursor:    1,
		curr_byte: 0b00000000,
	}

	if !areWritersEqual(w, expected_writer) || tree_encoded_size != uint32(MOCK_TREE_ENCODED_SIZE) {
		t.Fatalf(`Test Write Tree Failed.
		writer: %+v,
		expected: %+v`, w, expected_writer)
	}
}
