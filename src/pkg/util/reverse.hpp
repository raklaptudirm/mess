#ifndef UTIL_REVERSE
#define UTIL_REVERSE

#include "types.hpp"

// Lookup table with precomputed reverses for each byte value.
constexpr uint64 BitReverseTable256[256] = {
    #define R2(n)   (n),   (n + 2*64),   (n + 1*64),   (n + 3*64)
    #define R4(n) R2(n), R2(n + 2*16), R2(n + 1*16), R2(n + 3*16)
    #define R6(n) R4(n), R4(n + 2*4 ), R4(n + 1*4 ), R4(n + 3*4 )
    R6(0), R6(2), R6(1), R6(3)
};

// reverse reverses the bits of the given uint64 number.
constexpr inline uint64 reverse(uint64 n) {
    return (BitReverseTable256[(n >>  0) & 0xff] << 56) |
           (BitReverseTable256[(n >>  8) & 0xff] << 48) |
           (BitReverseTable256[(n >> 16) & 0xff] << 40) |
           (BitReverseTable256[(n >> 24) & 0xff] << 32) |
           (BitReverseTable256[(n >> 32) & 0xff] << 24) |
           (BitReverseTable256[(n >> 40) & 0xff] << 16) |
           (BitReverseTable256[(n >> 48) & 0xff] <<  8) |
           (BitReverseTable256[(n >> 56) & 0xff] <<  0);
}

#endif
