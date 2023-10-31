#include <iostream>

#include "chess/board.hpp"

#include "types/types.hpp"

// NOLINTNEXTLINE chrono header unnecessary in Darwin.
#include <chrono>

using namespace Chess;

int main(int argc, char const *argv[]) {
    const auto fen = argc >= 2 ? argv[1] : "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1";
    const auto dep = argc >= 3 ? argv[2] : "7";

    Board board = Board(fen);

    const auto start = std::chrono::steady_clock::now();
    const auto nodes = board.Perft<true, true>((int8)std::atoi(dep));
    const auto end = std::chrono::steady_clock::now();

    const std::chrono::duration<float64> delta = end - start;
    std::cout << "nodes " << nodes << " nps " << (uint64)((nodes / delta.count()) / 1'000'000) << std::endl;
}
