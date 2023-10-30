#include <iostream>

#include "chess/move.hpp"
#include "chess/square.hpp"
#include "chess/board.hpp"

#include "types/types.hpp"

// NOLINTNEXTLINE chrono header unnecessary in Darwin.
#include <chrono>

using namespace Chess;

template <bool BULK_COUNT, bool SPLIT_MOVES>
// NOLINTNEXTLINE(misc-no-recursion)
uint64 perft(Board& board, int8 depth) {
    // Return 1 for current node at depth 0.
    if (depth <= 0)
        return 1;

    // Generate legal move-list.
    const auto moves = board.GenerateMoves<true, true>();

    // When bulk counting is enabled, return the length of
    // the legal move-list when depth is one. This saves a
    // lot of time cause it saves make moves and recursion.
    if (BULK_COUNT && !SPLIT_MOVES && depth == 1)
        return moves.Length();

    // Variable to cumulate node count in.
    uint64 nodes = 0;

    // Recursively call perft for child nodes.
    for (const auto move : moves) {
        board.MakeMove(move);
        const uint64 delta = perft<BULK_COUNT, false>(board, depth - 1);
        board.UndoMove();

        nodes += delta;

        // If split moves is enabled, display each child move's
        // contribution to the node count separately.
        if (SPLIT_MOVES)
            std::cout << move << ": " << delta << std::endl;

    }

    // Return cumulative node count.
    return nodes;
}

int main(int argc, char const *argv[]) {
    const auto fen = argc >= 2 ? argv[1] : "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1";
    const auto dep = argc >= 3 ? argv[2] : "7";

    Board board = Board(fen);

    const auto start = std::chrono::steady_clock::now();
    const auto nodes = perft<true, true>(board, (int8)std::atoi(dep));
    const auto end = std::chrono::steady_clock::now();

    const std::chrono::duration<float64> delta = end - start;
    std::cout << "nodes " << nodes << " nps " << (uint64)((nodes / delta.count()) / 1'000'000) << std::endl;
}
