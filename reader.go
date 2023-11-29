package main

import (
	"fmt"
	"io"
)

type Reader struct {
	buf                        []byte
	buf_len                    int
	idx                        int   // 0..len(buffer), index of current byte to read
	cursor                     uint8 // 0..7, index of current bit to read in byte buffer[index]
	second_last_byte_read_size uint8
}

func GetReader(data []byte) *Reader {
	// in a correct huffman-encoded file, the last byte is an int between 0-7
	// for example if it's 3, then we only read the first 3 bits of the second-last byte
	// if it's 0, we read the whole 8 bits of the second-last byte
	// then we reach EOF

	last_byte := data[len(data)-1]
	return &Reader{
		buf:                        data,
		buf_len:                    len(data),
		second_last_byte_read_size: uint8(last_byte) % 8,
	}
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
	if r.idx == r.buf_len-2 {
		if r.cursor != 0 || r.second_last_byte_read_size != 0 {
			return b, fmt.Errorf("can't read byte at index %d. buffer [... %08b %08b], cursor %d, second_last_byte_read_size %d", r.idx, r.buf[r.idx], r.buf[r.idx+1], r.cursor, r.second_last_byte_read_size)
		}
	}

	if r.idx == r.buf_len-3 && r.cursor > r.second_last_byte_read_size && r.second_last_byte_read_size != 0 {
		return b, fmt.Errorf("can't read byte at index %d, cursor %d. buffer [... %08b %08b %08b], second_last_byte_read_size %d", r.idx, r.cursor, r.buf[r.idx], r.buf[r.idx+1], r.buf[r.idx+2], r.second_last_byte_read_size)
	}
	if r.cursor == 0 {
		b = r.buf[r.idx]
	} else {
		x, y := r.buf[r.idx], r.buf[r.idx+1]
		var x_mask byte = 0xff >> r.cursor
		var y_mask byte = 0xff << (8 - r.cursor)

		lastOfx := (x & x_mask) << r.cursor
		firstOfy := (y & y_mask) >> (8 - r.cursor)
		b = lastOfx | firstOfy
	}

	r.idx++
	return b, nil
}

func (r *Reader) ReadTree(read_start_index int, tree_size int) (*Node, error) {

	read_current_index := 8*r.idx + int(r.cursor)

	if read_current_index-read_start_index >= tree_size {
		return nil, io.EOF
	}

	bit, err := r.ReadBit()
	if err != nil {
		return nil, err
	}

	if bit == 1 {
		ch, read_byte_err := r.ReadByte()
		if read_byte_err != nil {
			return nil, read_byte_err
		}
		return &Node{ch: ch}, nil
	}

	left_node, read_left_err := r.ReadTree(read_start_index, tree_size)
	if read_left_err != nil {
		if read_left_err == io.EOF {
			return nil, nil
		}
		return nil, read_left_err
	}
	right_node, read_right_err := r.ReadTree(read_start_index, tree_size)
	if read_right_err != nil {
		if read_left_err == io.EOF {
			return nil, nil
		}
		return nil, read_right_err
	}
	return &Node{
		Left:  left_node,
		Right: right_node,
	}, nil
}
