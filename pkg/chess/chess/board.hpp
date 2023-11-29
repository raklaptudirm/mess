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

#ifndef CHESS_BOARD
#define CHESS_BOARD

#include <array>
#include <cassert>
#include <utility>
#include <cstdint>

#include "fen.hpp"
#include "move.hpp"
#include "moves.hpp"
#include "square.hpp"
#include "zobrist.hpp"
#include "movegen.hpp"
#include "movelist.hpp"
#include "castling.hpp"
#include "bitboard.hpp"
#include "position.hpp"

namespace Chess {
    class Board {
    private:
        const Castling::Info castlingInfo = Castling::Info();

        // Position stack.
        uint16_t top = 0; // Current top in stack.
        std::array<Position, Move::MaxInGame> history; // Stack.

        // Game ply-count. May differ from top by a constant.
        const uint16_t initialPlys = 0;

        // Boolean representing if the Board is an FRC/DFRC Board.
        const bool frc;

        // push pushes a new Position into the Position stack.
        inline void push() {
            top++; // Move the top pointer higher.

            // Bounds check.
            assert(top < Move::MaxInGame);
        }

        // pop pops the top Position from the Position stack.
        inline void pop() {
            top--; // Move the top pointer lower.

            // Bounds check;
            assert(top >= 0);
        }

        // doCastling is a utility function used to make the final half of the castling move.
        // The function assumes that the castling king and rook have already been removed.
        static inline void doCastling(Position& position, Castling::Side side) {
            // Get the final squares for the particular castling dimension.
            const auto dim = Castling::Dimension(position.SideToMove, side);
            const auto ends = Castling::EndSquares(dim);

            // Insert the castling king and rook into their final squares.
            position.Insert(ends.first,  Piece::King + position.SideToMove);
            position.Insert(ends.second, Piece::Rook + position.SideToMove);
        }

    public:
        // Constructor of board using a fen string.
        explicit Board(const FEN& fen) : castlingInfo(fen.CastlingInfo), initialPlys(fen.PlysCount), frc(fen.FRC) {
            // Create a new position from the fen and store it
            // at the top in the Position stack.
            history[top] = Chess::Position(fen);
        }

        // Position returns a reference to the current Board Position.
        [[nodiscard]] const Position& Position() const {
            return history[top];
        }

        // PlyCount returns the number of plys in the current game.
        [[maybe_unused]] [[nodiscard]] uint16_t PlyCount() const {
            // Number of plys is equal to initial number of plys at
            // root (which may be non-zero for non-startpos positions)
            // + the number of plys since the root (top).
            return initialPlys + top;
        }

        // MakeMove makes the given chess move on the Board. It does not
        // check the legality of the provided move and assumes that it
        // is legal, therefore making it the responsibility of the caller.
        void MakeMove(Move move) {
            push(); // Push a new position into the Position stack.

            // Copy the last position into the current position
            // so that we can make the move without editing.
            history[top] = history[top - 1];

            // Create a reference of the top position, so we can
            // reference and edit it easily without indexing history.
            auto& position = history[top];

            // Variables for the source and target squares of the move.
            const auto source = move.Source();
            const auto target = move.Target();

            // Variable for the move flag which stores metadata.
            const auto flag = move.Flag();

            // Variables for the pieces at the source and target squares
            // prior to making the chess move on the Board.
            const auto sourcePiece = position[source];
            const auto targetPiece = position[target];

            // If the target square is not empty, the move is a capture.
            // The function assumes that the provided move is legal, and
            // therefore the piece at the target square being a friendly
            // piece is impossible.
            const bool isCapture = targetPiece != ColoredPiece::None;

            // UP represents the up direction for the current side to move.
            const auto UP = Directions::Up(position.SideToMove);

            // Increase the draw clock. Any reset is done later in the code.
            position.DrawClock++;

            // Clear the en-passant target square, if any.
            if (position.EpTarget != Square::None) {
                position.Hash -= Keys::EnPassantTarget(position.EpTarget);
                position.EpTarget = Square::None;
            }

            // Determine the change in castling rights, if any.
            const Castling::Rights change = castlingInfo.Mask(source) + castlingInfo.Mask(target);
            position.Hash -= Keys::CastlingRights(change & position.Rights);
            position.Rights -= change;

            // Remove the moving piece.
            position.Remove(source);

            // Reset the draw clock if any conditions are met.
            // Remove the captured piece if any.
            if (isCapture) {
                position.Remove(target);
                position.DrawClock = 0; // Reset on capture (irreversible move)
            } else if (sourcePiece.Piece() == Piece::Pawn) {
                position.DrawClock = 0; // Reset on pawn move (irreversible move)
            }

            switch (static_cast<uint8_t>(flag)) {
                case MoveFlag::Normal: {
                    // Normal move, insert the moving piece to the target.
                    position.Insert(target, sourcePiece);
                    break;
                }

                case MoveFlag::DoublePush: {
                    position.Insert(target, sourcePiece);

                    // Pawn double push, set the en-passant square if there are
                    // enemy pawns which can capture en-passant on the next move.
                    const auto newEPTarget = source >> UP;
                    const auto attackers = position[Piece::Pawn] & position[!position.SideToMove];
                    if (!MoveTable::Pawn(position.SideToMove, newEPTarget).IsDisjoint(attackers)) {
                        position.EpTarget = newEPTarget;
                        position.Hash += Keys::EnPassantTarget(position.EpTarget);
                    }
                    break;
                }

                // Castling move, details handled by doCastling function.
                case MoveFlag::CastleHSide: doCastling(position, Castling::Side::H); break;
                case MoveFlag::CastleASide: doCastling(position, Castling::Side::A); break;

                case MoveFlag::EnPassant: {
                    position.Insert(target, sourcePiece);

                    // En Passant capture, remove the correct pawn.
                    position.Remove(target >> -UP);
                    break;
                }

                // Promotion move, insert the promoted piece to the target.
                case MoveFlag::QPromotion: position.Insert(target, Piece::Queen  + position.SideToMove); break;
                case MoveFlag::NPromotion: position.Insert(target, Piece::Knight + position.SideToMove); break;
                case MoveFlag::BPromotion: position.Insert(target, Piece::Bishop + position.SideToMove); break;
                case MoveFlag::RPromotion: position.Insert(target, Piece::Rook   + position.SideToMove); break;

                // All flags are handled above, unreachable code.
                default: assert(false);
            }

            // Switch side to move.
            position.SideToMove = !position.SideToMove;
            position.Hash += Keys::SideToMove;

            // Generate checker BitBoard.
            position.GenerateCheckers();

            // Ensure the incremental hash is equal to the correct hash.
            assert(position.Hash == Position::ZobristHash(position));
        }

        // UndoMove undoes the last chess move made on the Board.
        void UndoMove() {
            // Undoing a move is just popping the top Position
            // from the stack, making the last Position the new top.
            pop();
        }

        // GenerateMoves generates the legal moves in the current position which
        // follow the provided move-generation options, and returns a MoveList.
        template<bool QUIET, bool NOISY>
        [[nodiscard]] MoveList GenerateMoves() const {
            return Position().SideToMove == Color::White ?
                    Moves::Generate<Color::White, QUIET, NOISY>(Position(), castlingInfo) :
                    Moves::Generate<Color::Black, QUIET, NOISY>(Position(), castlingInfo);
        }

        // ToString converts the target Board into its string representation.
        [[nodiscard]] constexpr std::string ToString() const {
            return Position().ToString();
        }

        // ToString converts the given move to its string representation,
        // using the correct representation for standard castling moves.
        [[nodiscard]] constexpr std::string ToString(Move move) const {
            if (!frc) {
                // Use the king to king target expression
                // for castling in non-frc boards.
                switch (static_cast<uint8_t>(move.Flag())) {
                    case MoveFlag::CastleASide:
                        return Move(
                                move.Source(),
                                Castling::EndSquares(
                                        Castling::Dimension(
                                                Position().SideToMove,
                                                Castling::Side::A
                                        )
                                ).first,
                                MoveFlag::CastleASide
                        ).ToString();
                    case MoveFlag::CastleHSide:
                        return Move(
                                move.Source(),
                                Castling::EndSquares(
                                        Castling::Dimension(
                                                Position().SideToMove,
                                                Castling::Side::H
                                        )
                                ).first,
                                MoveFlag::CastleHSide
                        ).ToString();
                }
            }

            // Use the internal representation in all other cases.
            return move.ToString();
        }

        // Perft implements the perft function, which counts the number
        // of nodes at a given depth from the given position. BULK_COUNT
        // enables bulk counting which makes perft much faster, but is
        // unusable in standard search. SPLIT_MOVES gives a breakdown of
        // the nodes contributed by each move from the root.
        template <bool BULK_COUNT, bool SPLIT_MOVES>
        [[maybe_unused]] int64_t Perft(int32_t depth) {
            return perft<BULK_COUNT, SPLIT_MOVES>(*this, depth);
        }

    private:
        template <bool BULK_COUNT, bool SPLIT_MOVES>
        // NOLINTNEXTLINE(misc-no-recursion)
        static int64_t perft(Board& board, int32_t depth) {
            // Return 1 for current node at depth 0.
            if (depth <= 0)
                return 1;

            // Generate legal move-list.
            const auto moves = board.GenerateMoves<true, true>();

            // When bulk counting is enabled, return the length of
            // the legal move-list when depth is one. This saves a
            // lot of time cause it saves make moves and recursion.
            if (BULK_COUNT && !SPLIT_MOVES && depth == 1)
                return static_cast<int64_t>(moves.Length());

            // Variable to cumulate node count in.
            int64_t nodes = 0;

            // Recursively call perft for child nodes.
            for (const auto move : moves) {
                board.MakeMove(move);
                const int64_t delta = perft<BULK_COUNT, false>(board, depth - 1);
                board.UndoMove();

                nodes += delta;

                // If split moves is enabled, display each child move's
                // contribution to the node count separately.
                if (SPLIT_MOVES)
                    std::cout << board.ToString(move) << ": " << delta << std::endl;

            }

            // Return cumulative node count.
            return nodes;
        }
    };

    constexpr inline std::ostream& operator<<(std::ostream& os, const Board& board) {
        os << board.ToString();
        return os;
    }
}

#endif