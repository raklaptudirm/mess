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

#include "fen.hpp"
#include "move.hpp"
#include "moves.hpp"
#include "square.hpp"
#include "zobrist.hpp"
#include "castling.hpp"
#include "bitboard.hpp"
#include "position.hpp"

namespace Chess {
    class Board {
    private:
        Castling::Info castlingInfo = Castling::Info();

        // Position stack.
        uint16 top = 0; // Current top in stack.
        std::array<Position, 512> history = {}; // Stack.

        // Game ply-count. May differ from top by a constant.
        uint16 plys = 0;

        // push pushes a new Position into the Position stack.
        void push() {
            top++;
            plys++;
        }

        // pop pops the top Position from the Position stack.
        void pop() {
            top--;
            plys--;
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
        explicit Board(const std::string& fenString) {
            // Parse the provided fen string.
            const FEN fen = FEN(fenString);

            // Copy the fields relevant to Board.
            plys         = fen.PlysCount;
            castlingInfo = fen.CastlingInfo;

            // Create a new position from the fen and store it
            // at the top in the Position stack.
            history[top] = Chess::Position(fen);
        }

        // Position returns a reference to the current Board Position.
        [[nodiscard]] const Position& Position() const {
            return history[top];
        }

        // TODO: make this private
        [[nodiscard]] const Castling::Info& CastlingInfo() const {
            return castlingInfo;
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
            // TODO: add castling rights to zobrist hash
            const Castling::Rights change = castlingInfo.Mask(source) + castlingInfo.Mask(target);
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

            switch (flag) {
                case Move::Flag::Normal: {
                    // Normal move, insert the moving piece to the target.
                    position.Insert(target, sourcePiece);
                    break;
                }

                case Move::Flag::DoublePush: {
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
                case Move::Flag::CastleHSide: doCastling(position, Castling::Side::H); break;
                case Move::Flag::CastleASide: doCastling(position, Castling::Side::A); break;

                case Move::Flag::EnPassant: {
                    position.Insert(target, sourcePiece);

                    // En Passant capture, remove the correct pawn.
                    position.Remove(target >> -UP);
                    break;
                }

                // Promotion move, insert the promoted piece to the target.
                case Move::Flag::QPromotion: position.Insert(target, Piece::Queen  + position.SideToMove); break;
                case Move::Flag::NPromotion: position.Insert(target, Piece::Knight + position.SideToMove); break;
                case Move::Flag::BPromotion: position.Insert(target, Piece::Bishop + position.SideToMove); break;
                case Move::Flag::RPromotion: position.Insert(target, Piece::Rook   + position.SideToMove); break;

                // All flags are handled above, unreachable code.
                default: assert(false);
            }

            // Increase the number of game plys.
            plys++;

            // Switch side to move.
            position.SideToMove = !position.SideToMove;
            position.Hash += Keys::SideToMove;

            // Generate checker BitBoard.
            position.GenerateCheckers();
        }

        // UndoMove undoes the last chess move made on the Board.
        void UndoMove() {
            // Undoing a move is just popping the top Position
            // from the stack, making the last Position the new top.
            pop();
        }

        // ToString converts the target Board into its string representation.
        [[nodiscard]] constexpr std::string ToString() const {
            return Position().ToString();
        }
    };

    constexpr inline std::ostream& operator<<(std::ostream& os, const Board& board) {
        os << board.ToString();
        return os;
    }
}

#endif