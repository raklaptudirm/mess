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

// xorshift64star Pseudo-Random Number Generator
// This struct is based on original code written and dedicated
// to the public domain by Sebastiano Vigna (2014).
// It has the following characteristics:
//
//   - Outputs 64-bit numbers
//   - Passes Dieharder and SmallCrush test batteries
//   - Does not require warm-up, no zeroland to escape
//   - Internal state is a single 64-bit integer
//   - Period is 2^64 - 1
//   - Speed: 1.60 ns/call (Core i7 @3.40GHz)
//
// For further analysis see
//
//	<http://vigna.di.unimi.it/ftp/papers/xorshift.pdf>
type PRNG struct {
	seed uint64
}

// Seed seeds the pseudo-random number generator with the given uint.
func (p *PRNG) Seed(s uint64) {
	p.seed = s
}

// Uint64 generates a new pseudo-random uint64.
func (p *PRNG) Uint64() uint64 {
	// linear feedback shifts
	p.seed ^= p.seed >> 12
	p.seed ^= p.seed << 25
	p.seed ^= p.seed >> 27

	// scramble result with non-linear function
	return p.seed * 2685821657736338717
}

// SparseUint64 generates a pseudo-random sparse uint64, i.e, a number
// with very few set bits. This is useful in magic table generation.
func (p *PRNG) SparseUint64() uint64 {
	// bitwise and three pseudo-random uint64s together
	// only the bits set in all three are set in the result

	//nolint:staticcheck // Uint64 is an impure function
	return p.Uint64() & p.Uint64() & p.Uint64()
}
