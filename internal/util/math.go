package util

// Type integer represents every value that can be represented as an integer.
type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Max returns the larger value between the integers a and b.
func Max[T integer](a, b T) T {
	if a > b {
		return a
	}

	return b
}

// Min returns the smaller value between the integers a and b.
func Min[T integer](a, b T) T {
	if a < b {
		return a
	}

	return b
}

// Abs returns the absolute value of the integer x.
func Abs[T integer](x T) T {
	if x < 0 {
		return -x
	}

	return x
}
