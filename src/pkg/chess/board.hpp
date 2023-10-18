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
        public:
            Castling::Info CastlingInfo = Castling::Info();

        private:

            // 8x8 Mailbox position representation.
            uint16 offset = 0;
            std::array<Position, 512> history = {};

            uint16 plysCount = 0;

            constexpr inline Board() = default;

        public:
            explicit Board(const std::string& fenString) {
                const FEN fen = FEN(fenString);

                plysCount = fen.PlysCount;
                CastlingInfo = fen.CastlingInfo;

                offset = 0;
                history[offset] = Chess::Position(fen);
            }

            [[nodiscard]] const Position& Position() const {
                return history[offset];
            }

            void MakeMove(Move move) {
                history[offset + 1] = history[offset];
                offset++;
                auto& position = history[offset];
                plysCount++;

                const auto source = move.Source();
                const auto target = move.Target();

                const auto flag = move.Flag();

                const auto sourcePiece = position[source];
                const auto targetPiece = position[target];

                const bool isCapture = targetPiece != ColoredPiece::None;

                const auto UP = Directions::Up(position.SideToMove);

                position.DrawClock++;

                if (position.EpTarget != Square::None) {
                    position.Hash -= Keys::EnPassantTarget(position.EpTarget);
                    position.EpTarget = Square::None;
                }

                const Castling::Rights change = CastlingInfo.Mask(source) + CastlingInfo.Mask(target);
                position.Rights -= change;

                position.Remove(source);
                if (isCapture) {
                    position.Remove(target);
                    position.DrawClock = 0;
                } else if (sourcePiece.Piece() == Piece::Pawn) {
                    position.DrawClock = 0;
                }

                switch (flag) {
                    case Move::Flag::Normal: {
                        position.Insert(target, sourcePiece);
                        break;
                    }

                    case Move::Flag::DoublePush: {
                        position.Insert(target, sourcePiece);

                        const auto newEPTarget = source >> UP;
                        const auto attackers = position[Piece::Pawn] & position[!position.SideToMove];
                        if (!MoveTable::Pawn(position.SideToMove, newEPTarget).IsDisjoint(attackers)) {
                            position.EpTarget = newEPTarget;
                            position.Hash += Keys::EnPassantTarget(position.EpTarget);
                        }
                        break;
                    }

                    case Move::Flag::CastleHSide: {
                        const auto dim = Castling::Dimension(position.SideToMove, Castling::Side::H);
                        const auto ends = Castling::EndSquares(dim);
                        position.Insert(ends.first,  Piece::King + position.SideToMove);
                        position.Insert(ends.second, Piece::Rook + position.SideToMove);
                        break;
                    }
                    case Move::Flag::CastleASide: {
                        const auto dim = Castling::Dimension(position.SideToMove, Castling::Side::A);
                        const auto ends = Castling::EndSquares(dim);
                        position.Insert(ends.first,  Piece::King + position.SideToMove);
                        position.Insert(ends.second, Piece::Rook + position.SideToMove);
                        break;
                    }

                    case Move::Flag::EnPassant: {
                        position.Insert(target, sourcePiece);
                        position.Remove(target >> -UP);
                        break;
                    }

                    case Move::Flag::QPromotion: position.Insert(target, Piece::Queen  + position.SideToMove); break;
                    case Move::Flag::NPromotion: position.Insert(target, Piece::Knight + position.SideToMove); break;
                    case Move::Flag::BPromotion: position.Insert(target, Piece::Bishop + position.SideToMove); break;
                    case Move::Flag::RPromotion: position.Insert(target, Piece::Rook   + position.SideToMove); break;

                    default: assert(false);
                }

                plysCount++;

                position.SideToMove = !position.SideToMove;
                position.Hash += Keys::SideToMove;

                position.GenerateCheckers();
            }

            void UndoMove() {
                offset--;
                plysCount--;
            }

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