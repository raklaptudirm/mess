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

namespace Chess {
    class Board {
        public:
            struct State {
                // Information about Child Node.
                Move         PlayedMove;    // Move Played on this Position.
                ColoredPiece CapturedPiece; // Piece captured by the Move.

                // Irreversible Information.
                Castling::Rights CurrentRights; // Current Castling Rights.
                Square CurrentEPTarget;  // Current En Passant Target.
                uint8  CurrentDrawClock; // Current Draw Clock State.

                // Reversible but complicated to
                // recalculate from scratch.
                Hash CurrentHash; // Current Position Hash.
            };

        private:

            // 8x8 Mailbox board representation.
            std::array<ColoredPiece, Square::N> mailbox = {};

            // BitBoard Board representation.
            // x6 piece bbs, x2 color bbs.
            std::array<BitBoard, Piece::N> pieceBBs = {};
            std::array<BitBoard, Color::N> colorBBs = {};

            Color sideToMove = Color ();
            Square epTarget  = Square();

            Castling::Info castling = Castling::Info();

            uint16 plysCount = 0;
            uint8  drawClock = 0;

            Hash hash = Hash();

            BitBoard checkers = BitBoard();
            uint8    checkNum = 0;

            bool frc = false;

            std::array<State, Move::MaxInGame*2> history = {};

            constexpr inline Board() = default;

            constexpr inline void insert(Square square, ColoredPiece piece) {
                assert(square != Square::None && piece != ColoredPiece::None);

                mailbox[static_cast<uint8>(square)] = piece;

                pieceBBs[static_cast<uint8>(piece.Piece())].Flip(square);
                colorBBs[static_cast<uint8>(piece.Color())].Flip(square);

                hash += Keys::PieceOnSquare(piece, square);
            }

            constexpr inline void remove(Square square) {
                const ColoredPiece piece = mailbox[static_cast<uint8>(square)];

                assert(square != Square::None);
                assert(piece != ColoredPiece::None);

                mailbox[static_cast<uint8>(square)] = ColoredPiece::None;

                pieceBBs[static_cast<uint8>(piece.Piece())].Flip(square);
                colorBBs[static_cast<uint8>(piece.Color())].Flip(square);

                hash -= Keys::PieceOnSquare(piece, square);
            }

        public:

            /********************************************
             * Getter Functions for readonly properties *
             ********************************************/
            [[nodiscard]] constexpr inline Color    SideToMove()   const { return sideToMove; } // Side To Move.
            [[nodiscard]] constexpr inline Square   EPTarget()     const { return epTarget;   } // En Passant Target.
            [[nodiscard]] constexpr inline uint16   PlyCount()     const { return plysCount;  } // No of Plys since Root.
            [[nodiscard]] constexpr inline uint16   DrawClock()    const { return drawClock;  } // 50 move rule Draw Clock.
            [[nodiscard]] constexpr inline bool     FisherRandom() const { return frc;        } // Is a Fisher Random Position.
            [[nodiscard]] constexpr inline BitBoard Checkers()     const { return checkers;   }
            [[nodiscard]] constexpr inline uint8    CheckNum()     const { return checkNum;   }
            [[nodiscard]] constexpr inline Castling::Info* Castling() { return &castling; }


            // Indexing with a Piece return's that Piece's BitBoard.
            constexpr inline BitBoard operator [](const Piece piece) const {
                return pieceBBs[static_cast<uint8_t>(piece)];
            }

            // Indexing with a Color return's that Color's BitBoard.
            constexpr inline BitBoard operator [](const Color color) const {
                return colorBBs[static_cast<uint8_t>(color)];
            }

            /******************************
             * Game Termination Detection *
             ******************************/

            [[nodiscard]] constexpr inline bool Mated() const {
                return checkers.Some();
            }

            [[nodiscard]] constexpr inline bool Draw() const {
                return DrawBy50Move() || DrawByRepetition();
            }

            [[nodiscard]] constexpr inline bool DrawBy50Move() const {
                return drawClock >= 100 && !Mated();
            }

            [[nodiscard]] constexpr inline bool DrawByRepetition() const {
                return false;
            }

            [[nodiscard]] constexpr inline bool DrawByInsufficientMaterial() const {
                if ((*this)[Piece::Pawn].Some() || (*this)[Piece::Rook].Some() || (*this)[Piece::Queen].Some())
                    return false;

                return true;
            }

            /****************************
             * Public Utility Functions *
             ****************************/

            template <Color BY>
            constexpr inline bool Checked() {
                return Attacked<!BY>(((*this)[Piece::King] & (*this)[BY]).LSB());
            }

            template <Color BY>
            [[nodiscard]] constexpr inline bool Attacked(const Square square, const BitBoard blockers) const {
                const BitBoard attackers = (*this)[BY];

                // Check for pawn attackers.
                const BitBoard attackingPawns = (*this)[Piece::Pawn] & attackers;
                if (!attackingPawns.IsDisjoint(MoveTable::Pawn<!BY>(square))) return true;

                // Check for knight attackers.
                const BitBoard attackingKnights = (*this)[Piece::Knight] & attackers;
                if (!attackingKnights.IsDisjoint(MoveTable::Knight(square))) return true;

                const BitBoard attackingQueens = (*this)[Piece::Queen];

                // Check for bishop type attackers.
                const BitBoard attackingBishops = ((*this)[Piece::Bishop] | attackingQueens) & attackers;
                if (!attackingBishops.IsDisjoint(MoveTable::Bishop(square, blockers))) return true;

                // Check for rook type attackers.
                const BitBoard attackingRooks = ((*this)[Piece::Rook] | attackingQueens) & attackers;
                if (!attackingRooks.IsDisjoint(MoveTable::Rook(square, blockers))) return true;

                // Check for a king attacker.
                const BitBoard attackingKing = (*this)[Piece::King] & attackers;
                if (!attackingKing.IsDisjoint(MoveTable::King(square))) return true;

                // No attackers found.
                return false;
            }

            template <Color BY>
            [[nodiscard]] constexpr inline bool Attacked(const Square square) const {
                return Attacked<BY>(square, (*this)[BY] | (*this)[!BY]);
            }

            template<Color BY>
            [[nodiscard]] inline bool Attacked(const BitBoard targets, const BitBoard blockers) const {
                for (const auto target : targets)
                    if (Attacked<BY>(target, blockers)) return true;
                return false;
            }

            template<Color BY>
            [[nodiscard]] inline bool Attacked(const BitBoard targets) const {
                return Attacked<BY>(targets, (*this)[BY] | (*this)[!BY]);
            }

            explicit Board(const std::string& fenString) {
                const FEN fen = FEN(fenString);

                assert(sideToMove == Color ::None); sideToMove = fen.SideToMove;
                assert(epTarget   == Square::None); epTarget   = fen.EPTarget;

                assert(plysCount == 0); plysCount  = fen.PlysCount;
                assert(drawClock == 0); drawClock  = fen.DrawClock;

                castling = fen.CastlingInfo;

                for (uint8 sq = 0; sq < Square::N; sq++)
                    if (fen.Mailbox[sq] != ColoredPiece::None)
                        insert(Square(sq), fen.Mailbox[sq]);

                generateCheckers();
            }

            constexpr inline ColoredPiece operator[](const Square sq) const {
                return mailbox[static_cast<uint8>(sq)];
            }

            void MakeMove(Move move) {
                const auto source = move.Source();
                const auto target = move.Target();

                const auto flag = move.Flag();

                const auto sourcePiece = mailbox[static_cast<uint8>(source)];
                const auto targetPiece = mailbox[static_cast<uint8>(target)];

                const bool isCapture = targetPiece != ColoredPiece::None;

                const auto UP = Directions::Up(sideToMove);

                history[plysCount].PlayedMove    = move;
                history[plysCount].CapturedPiece = targetPiece;

                if (history[plysCount].CurrentHash != hash) {
                    history[plysCount].CurrentHash      = hash;
                    history[plysCount].CurrentRights    = castling.Rights;
                    history[plysCount].CurrentEPTarget  = epTarget;
                    history[plysCount].CurrentDrawClock = drawClock;
                }

                if (isCapture || sourcePiece.Piece() == Piece::Pawn)
                    drawClock = 0;
                else drawClock++;

                if (epTarget != Square::None) {
                    hash -= Keys::EnPassantTarget(epTarget);
                    epTarget = Square::None;
                }

                const Castling::Rights change = castling.Mask(source) + castling.Mask(target);
                if (change != Castling::None) {
                    castling.Rights -= change;
                }

                remove(source);
                if (isCapture) remove(target);

                switch (flag) {
                    case Move::Flag::Normal: {
                        insert(target, sourcePiece);
                        break;
                    }

                    case Move::Flag::DoublePush: {
                        insert(target, sourcePiece);

                        const auto newEPTarget = source >> UP;
                        const auto attackers = (*this)[Piece::Pawn] & (*this)[!sideToMove];
                        if (!MoveTable::Pawn(sideToMove, newEPTarget).IsDisjoint(attackers)) {
                            epTarget = newEPTarget;
                            hash += Keys::EnPassantTarget(epTarget);
                        }
                        break;
                    }

                    case Move::Flag::CastleHSide: {
                        const Castling::EndSquares ends = Castling::HEndSquares(sideToMove);
                        insert(ends.King, ColoredPiece(Piece::King, sideToMove));
                        insert(ends.Rook, ColoredPiece(Piece::Rook, sideToMove));
                        break;
                    }
                    case Move::Flag::CastleASide: {
                        const Castling::EndSquares ends = Castling::AEndSquares(sideToMove);
                        insert(ends.King, ColoredPiece(Piece::King, sideToMove));
                        insert(ends.Rook, ColoredPiece(Piece::Rook, sideToMove));
                        break;
                    }

                    case Move::Flag::EnPassant: {
                        insert(target, sourcePiece);
                        remove(target >> -UP);
                        break;
                    }

                    case Move::Flag::QPromotion: insert(target, ColoredPiece(Piece::Queen,  sideToMove)); break;
                    case Move::Flag::NPromotion: insert(target, ColoredPiece(Piece::Knight, sideToMove)); break;
                    case Move::Flag::BPromotion: insert(target, ColoredPiece(Piece::Bishop, sideToMove)); break;
                    case Move::Flag::RPromotion: insert(target, ColoredPiece(Piece::Rook,   sideToMove)); break;

                    default: assert(false);
                }

                plysCount++;

                sideToMove = !sideToMove;
                hash += Keys::SideToMove;

                generateCheckers();
            }

            void UndoMove() {
                plysCount--;
                sideToMove = !sideToMove;

                const auto move = history[plysCount].PlayedMove;

                const auto source = move.Source();
                const auto target = move.Target();

                const auto flag = move.Flag();

                const auto targetPiece = mailbox[static_cast<uint8>(target)];
                const auto sourcePiece = Move::Flag::IsPromotion(flag) ? ColoredPiece(Piece::Pawn, sideToMove)
                                                                       : targetPiece;

                switch (flag) {
                    case Move::Flag::CastleHSide: {
                        const Castling::EndSquares ends = Castling::HEndSquares(sideToMove);
                        if (source != ends.King) {
                            insert(source, mailbox[static_cast<uint8>(ends.King)]);
                            remove(ends.King);
                        }

                        if (target != ends.Rook) {
                            insert(target, mailbox[static_cast<uint8>(ends.Rook)]);
                            remove(ends.Rook);
                        }

                        break;
                    }

                    case Move::Flag::CastleASide: {
                        const Castling::EndSquares ends = Castling::AEndSquares(sideToMove);

                        if (source != ends.King) {
                            insert(source, mailbox[static_cast<uint8>(ends.King)]);
                            remove(ends.King);
                        }

                        if (target != ends.Rook) {
                            insert(target, mailbox[static_cast<uint8>(ends.Rook)]);
                            remove(ends.Rook);
                        }

                        break;
                    }

                    case Move::Flag::EnPassant: {
                        remove(target);
                        insert(target >> Directions::Down(sideToMove), ColoredPiece(Piece::Pawn, !sideToMove));
                        insert(source, sourcePiece);
                        break;
                    }

                    default:
                        remove(target);

                        if (history[plysCount].CapturedPiece != ColoredPiece::None)
                            insert(target, history[plysCount].CapturedPiece);

                        insert(source, sourcePiece);
                }

                epTarget  = history[plysCount].CurrentEPTarget;
                drawClock = history[plysCount].CurrentDrawClock;

                castling.Rights = history[plysCount].CurrentRights;

                hash = history[plysCount].CurrentHash;

                generateCheckers();
            }

            [[nodiscard]] constexpr std::string ToString() const {
                std::string board = "+---+---+---+---+---+---+---+---+\n";

                for (uint8 rank = 7; rank != 255; rank--) {
                    board += "| ";

                    for (uint8 file = 0; file < File::N; file++) {
                        board += mailbox[rank*8 + file].ToString() + " | ";
                    }

                    board += Rank(rank).ToString();
                        board += "\n+---+---+---+---+---+---+---+---+\n";
                }

                board += "  a   b   c   d   e   f   g   h\n";
                return board;
            }

            private:
                constexpr inline void generateCheckers() {
                    const auto friends = (*this)[ sideToMove];
                    const auto enemies = (*this)[!sideToMove];
                    const auto occupied = friends + enemies;

                    assert((*this)[Piece::King] != BitBoards::Empty);
                    const auto king = ((*this)[Piece::King] & friends).LSB();
                    assert(king != Square::None);

                    const auto p = (*this)[Piece::Pawn  ] & enemies;
                    const auto n = (*this)[Piece::Knight] & enemies;
                    const auto b = (*this)[Piece::Bishop] & enemies;
                    const auto r = (*this)[Piece::Rook  ] & enemies;
                    const auto q = (*this)[Piece::Queen ] & enemies;

                    const auto checkingP = p & MoveTable::Pawn(sideToMove, king);
                    const auto checkingN = n & MoveTable::Knight(king);
                    const auto checkingD = (b + q) & MoveTable::Bishop(king, occupied);
                    const auto checkingL = (r + q) & MoveTable::Rook  (king, occupied);

                    checkers = checkingP + checkingN + checkingD + checkingL;
                    checkNum = checkers.PopCount();
                }
    };

    constexpr inline std::ostream& operator<<(std::ostream& os, const Board& board) {
        os << board.ToString();
        return os;
    }
};

#endif