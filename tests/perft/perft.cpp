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

#ifndef TESTS_PERFT
#define TESTS_PERFT

#include <array>
#include <string>
#include <fstream>
#include <iostream>

#include "common.hpp"

#include "chess/board.hpp"
#include "strutil/strutil.h"
#include "catch2/catch_test_macros.hpp"

using namespace Chess;

namespace Perft {
    // TestCase represents a single PERFT test case which contains
    // a position and its PERFT results for depths 1-6.
    struct TestCase {
        std::string FEN;               // Position FEN Representation.
        std::array<int64, 7> Expected; // Expected PERFT Results for each depth.

        // TestCase creates a new TestCase from a string
        // in the format: fen { ;perft expected results }.
        // NOLINTNEXTLINE
        explicit TestCase(std::string caseStr) {
            strutil::trim(caseStr);
            auto fields = strutil::split(caseStr, ";");

            strutil::trim(fields[0]);
            FEN = fields[0];
            for (std::size_t i = 1; i < fields.size(); i++) {
                const auto perftResult = parsePerftResult(fields[i]);
                Expected[perftResult.first] = perftResult.second;
            }
        }

        // parsePerftResult parses a perft result case in the format:
        // ;D<depth> <expected perft result for depth in current position>
        static std::pair <int32, int64> parsePerftResult(std::string resultStr) {
            strutil::trim(resultStr);
            auto fields = strutil::split(resultStr, " ");

            // The call to string::substr removes the 'D' from "D<depth>" so that
            // we can parse "<depth>" as an integer without any errors.
            return {strutil::parse_string<int32>(fields[0].substr(1)),
                    strutil::parse_string<int64>(fields[1])};
        }
    };

    struct TestCases {
        // List of PERFT test cases.
        std::vector<TestCase> Cases = {};

        // TestCases reads a file containing PERFT test cases and
        // parses it into a usable TestCases object.
        explicit TestCases(const std::string& filename) {
            // Read the given file from the tests directory.
            std::ifstream file_in(TESTS_DIR + filename);

            // Check for any errors while reading the file.
            if (!file_in) {
                std::cout << "perft test: couldn't read " << filename << std::endl;
                return;
            }

            std::string line;

            // Parse each line as a separate TestCase.
            while (std::getline(file_in, line))
                Cases.emplace_back(line);
        }

        // Run the test cases stored in the current object.
        void Run(const int32 depth) const {
            const auto& tests = *this;
            int32 n = 1;
            for (const auto& test : tests.Cases) {
                // Print information in the format: [<index>/<total>] <fen>
                std::cout << "[" << std::setfill(' ') << std::setw(3) << n << "/"
                          << tests.Cases.size() << "] " << test.FEN << std::endl;
                n++; // Increase the index

                // Test the perft results from Chess::Board.
                auto chessboard = Chess::Board(test.FEN);
                CHECK(chessboard.Perft<true, false>(depth) == test.Expected[depth]);
            }
        }
    };
}

// PERFT_TEST creates a Catch2 PERFT test for the given variant and depth.
#define PERFT_TEST(type, depth)                                                      \
    TEST_CASE("perft: " type " depth "#depth, "[perft][" type "][depth "#depth"]") { \
        Perft::TestCases tests{"perft/" type ".epd"};                                \
        tests.Run(depth);                                                            \
    }

// PERFT_DEPTH_TEST creates Catch2 PERFT tests for the given depth and
// for all the supported variants, that is standard and chess 960.
#define PERFT_DEPTH_TEST(depth)   \
    PERFT_TEST("standard", depth) \
    PERFT_TEST("chess960", depth)

// PERFT tests for depths 1-6.
PERFT_DEPTH_TEST(1)
PERFT_DEPTH_TEST(2)
PERFT_DEPTH_TEST(3)
PERFT_DEPTH_TEST(4)
PERFT_DEPTH_TEST(5)
PERFT_DEPTH_TEST(6)

#endif //TESTS_PERFT
