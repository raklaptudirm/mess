// Copyright © 2023 Rak Laptudirm <rak@laptudirm.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

// Type number represents every value that can be represented as an number.
type number interface {
	integer | ~float32 | ~float64
}

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Max returns the larger value between the integers a and b.
func Max[T number](a, b T) T {
	if a > b {
		return a
	}

	return b
}

// Min returns the smaller value between the integers a and b.
func Min[T number](a, b T) T {
	if a < b {
		return a
	}

	return b
}

// Abs returns the absolute value of the integer x.
func Abs[T number](x T) T {
	if x < 0 {
		return -x
	}

	return x
}

// Clamp returns a value between min and max. If n is between min and max,
// n is returned, otherwise the closest limit is returned.
func Clamp[T number](n, min, max T) T {
	return Max(min, Min(n, max))
}

// Lerp does a linear interpolation between the given start and stop
// numbers by the ratio numerator/denominator(better precision for ints).
func Lerp[T number](start, stop, numerator, denominator T) T {
	return (stop*numerator + start*(denominator-numerator)) / denominator
}

// Btoi converts the given boolean value to a number. The number is
// equivalent to 1 if the value is true and 0 if it is false.
func Btoi[T number](condition bool) T {
	if condition {
		return 1
	}

	return 0
}
