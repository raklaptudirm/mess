#include <iostream>

#include "pkg/chess/move.hpp"
#include "pkg/chess/square.hpp"
#include "pkg/chess/board.hpp"
#include "pkg/chess/movegen.hpp"

#include <chrono>

using namespace Chess;

template <bool BULK_COUNT, bool SPLIT_MOVES>
uint64 perft(Board& board, int8 depth) {
    // Return 1 for current node at depth 0.
    if (depth <= 0)
        return 1;

    // When bulk counting is enabled, return the length of
    // the legal move-list when depth is one. This saves a
    // lot of time cause it saves make moves and recursion.
    if (BULK_COUNT && !SPLIT_MOVES && depth == 1)
        return Moves::Generate<true, true>(board).Length();

    // Generate legal move-list.
    const auto moves = Moves::Generate<true, true>(board);

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
    assert(argc == 3);
    Board board = Board(argv[1]);

    const auto start = std::chrono::steady_clock::now();
    const auto nodes = perft<true, true>(board, (int8)std::atoi(argv[2]));
    const auto end = std::chrono::steady_clock::now();

    const std::chrono::duration<float64> delta = end - start;
    std::cout << "nodes " << nodes << " nps " << (uint64)((nodes / delta.count()) / 1'000'000) << std::endl;
}
