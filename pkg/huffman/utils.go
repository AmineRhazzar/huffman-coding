package huffman
import (
	"encoding/binary"
	"bytes"
)

func Insert[T any](a *[]T, index int, value T) {
	if len((*a)) == index { // nil or empty slice or after last element
		(*a) = append((*a), value)
		return
	}
	(*a) = append((*a)[:index+1], (*a)[index:]...) // index < len(a)
	(*a)[index] = value
}

func ByteEncode(num uint32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}