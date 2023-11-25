package huffman

func Insert[T any](a *[]T, index int, value T) {
	if len((*a)) == index { // nil or empty slice or after last element
		(*a) = append((*a), value)
		return
	}
	(*a) = append((*a)[:index+1], (*a)[index:]...) // index < len(a)
	(*a)[index] = value
}
