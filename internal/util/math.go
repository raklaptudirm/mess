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
