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

#ifndef UTIL_TYPES
#define UTIL_TYPES

#include <cstdint>

/*****************
 * Integer Types *
 *****************/

using int8  [[maybe_unused]] = int8_t;  //  8-bit signed integer.
using int16 [[maybe_unused]] = int16_t; // 16-bit signed integer.
using int32 [[maybe_unused]] = int32_t; // 32-bit signed integer.
using int64 [[maybe_unused]] = int64_t; // 64-bit signed integer.

/**************************
 * Unsigned Integer Types *
 **************************/

using uint8  [[maybe_unused]] = uint8_t;  //  8-bit unsigned integer.
using uint16 [[maybe_unused]] = uint16_t; // 16-bit unsigned integer.
using uint32 [[maybe_unused]] = uint32_t; // 32-bit unsigned integer.
using uint64 [[maybe_unused]] = uint64_t; // 64-bit unsigned integer.

/************************
 * Floating Point Types *
 ************************/

using float32 [[maybe_unused]] = float;  // 32-bit floating point number.
using float64 [[maybe_unused]] = double; // 64-bit floating point number.

/***************
 * Alias Types *
 ***************/

using byte [[maybe_unused]] = uint8; // A single byte of memory. An alias of uint8.

#endif