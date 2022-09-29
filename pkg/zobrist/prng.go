package zobrist

// xorshift64star Pseudo-Random Number Generator
// This struct is based on original code written and dedicated
// to the public domain by Sebastiano Vigna (2014).
// It has the following characteristics:
//
//  -  Outputs 64-bit numbers
//  -  Passes Dieharder and SmallCrush test batteries
//  -  Does not require warm-up, no zeroland to escape
//  -  Internal state is a single 64-bit integer
//  -  Period is 2^64 - 1
//  -  Speed: 1.60 ns/call (Core i7 @3.40GHz)
//
// For further analysis see
//   <http://vigna.di.unimi.it/ftp/papers/xorshift.pdf>
type PRNG struct {
	seed uint64
}

func (p *PRNG) Seed(s uint64) {
	p.seed = s
}

func (p *PRNG) Uint64() uint64 {
	// linear feedback shifts
	p.seed ^= p.seed >> 12
	p.seed ^= p.seed << 25
	p.seed ^= p.seed >> 27

	// scramble result with non-linear function
	return p.seed * 2685821657736338717
}
