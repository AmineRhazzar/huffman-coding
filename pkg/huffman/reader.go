package huffman

import (
	"fmt"
	"io"
)

type Reader struct {
	buf                        []byte
	buf_len                    int
	idx                        int // 0..len(buffer), index of current byte to read
	cursor                     int // 0..7, index of current bit to read in byte buffer[index]
	second_last_byte_read_size int
	debug                      bool
}

func GetReader(data []byte, debug bool) *Reader {
	// in a correct huffman-encoded file, the last byte is an int between 0-7
	// for example if it's 3, then we only read the first 3 bits of the second-last byte
	// if it's 0, we read the whole 8 bits of the second-last byte
	// then we reach EOF

	last_byte := data[len(data)-1]
	return &Reader{
		buf:                        data,
		buf_len:                    len(data),
		second_last_byte_read_size: int(last_byte) % 8,
		debug:                      debug,
	}
}

func (r *Reader) Debug() {
	fmt.Printf("cursor: %d, index: %d, second_last_byte_read_size: %d\n", r.cursor, r.idx, r.second_last_byte_read_size)
}

func (r *Reader) ReadBit() (uint8, error) {
	if r.idx >= r.buf_len-1 {
		return 0, io.EOF
	}

	if r.idx == r.buf_len-2 &&
		r.second_last_byte_read_size != 0 &&
		r.cursor > r.second_last_byte_read_size-1 {
		return 0, io.EOF
	}

	curr_byte := r.buf[r.idx]
	bit := getBit(curr_byte, r.cursor)
	if r.cursor == 7 {
		r.cursor = 0
		r.idx++
	} else {
		r.cursor++
	}

	return bit, nil
}

func (r *Reader) ReadByte() (byte, error) {
	var b byte = 0
	// the last byte only has info on how many bits to read from the second last byte
	if r.idx >= r.buf_len-1 {
		return b, io.EOF
	}

	// normally this shouldn't happen in a correct huffman-encoded file
	if r.idx == r.buf_len-2 && r.second_last_byte_read_size != 0 {
		return b, fmt.Errorf("can't read byte at index %d. buffer [... %08b %08b], cursor %d, second_last_byte_read_size %d", r.idx, r.buf[r.idx], r.buf[r.idx+1], r.cursor, r.second_last_byte_read_size)
	}

	if r.cursor == 0 {
		b = r.buf[r.idx]
	} else {
		x, y := r.buf[r.idx], r.buf[r.idx+1]
		var x_mask byte = x << r.cursor
		var y_mask byte = y << (8 - r.cursor)
		b = x_mask | y_mask
	}

	r.idx++
	return b, nil
}

/**
1001001 01001100 01011001
*/
