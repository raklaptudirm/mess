package util

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func Max[T number](a, b T) T {
	if a > b {
		return a
	}

	return b
}

func Min[T number](a, b T) T {
	if a < b {
		return a
	}

	return b
}

func Abs[T number](x T) T {
	if x < 0 {
		return -x
	}

	return x
}
