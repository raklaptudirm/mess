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

#ifndef CHESS_POSITION
#define CHESS_POSITION

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
    struct Position {
        // 8x8 Mailbox position representation.
        std::array<ColoredPiece, Square::N> Mailbox = {};

        // BitBoard Board representation.
        // x6 piece bbs, x2 color bbs.
        std::array<BitBoard, Piece::N> PieceBBs = {};
        std::array<BitBoard, Color::N> ColorBBs = {};

        Hash Hash = Chess::Hash();
        BitBoard Checkers = BitBoard();

        Castling::Rights Rights = Castling::None;

        Color  SideToMove = Color ();
        Square EpTarget   = Square();

        uint8 DrawClock = 0;
        uint8 CheckNum  = 0;

        constexpr inline Position() = default;

        constexpr inline void Insert(Square square, ColoredPiece piece) {
            assert(square != Square::None && piece != ColoredPiece::None);

            Mailbox[static_cast<uint8>(square)] = piece;

            PieceBBs[static_cast<uint8>(piece.Piece())].Flip(square);
            ColorBBs[static_cast<uint8>(piece.Color())].Flip(square);

            Hash += Keys::PieceOnSquare(piece, square);
        }

        constexpr inline void Remove(Square square) {
            const ColoredPiece piece = Mailbox[static_cast<uint8>(square)];

            assert(square != Square::None);
            assert(piece != ColoredPiece::None);

            Mailbox[static_cast<uint8>(square)] = ColoredPiece::None;

            PieceBBs[static_cast<uint8>(piece.Piece())].Flip(square);
            ColorBBs[static_cast<uint8>(piece.Color())].Flip(square);

            Hash -= Keys::PieceOnSquare(piece, square);
        }

        // Indexing with a Piece return's that Piece's BitBoard.
        constexpr inline BitBoard operator [](const Piece piece) const {
            return PieceBBs[static_cast<uint8_t>(piece)];
        }

        // Indexing with a Color return's that Color's BitBoard.
        constexpr inline BitBoard operator [](const Color color) const {
            return ColorBBs[static_cast<uint8_t>(color)];
        }

        /******************************
         * Game Termination Detection *
         ******************************/

        [[nodiscard]] constexpr inline bool Mated() const {
            return Checkers.Some();
        }

        [[nodiscard]] constexpr inline bool Draw() const {
            return DrawBy50Move() || DrawByInsufficientMaterial();
        }

        [[nodiscard]] constexpr inline bool DrawBy50Move() const {
            return DrawClock >= 100 && !Mated();
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

        explicit Position(const std::string& fenString) {
            *this = Position(FEN(fenString));
        }

        explicit Position(const FEN& fen) {
            assert(SideToMove == Color ::None); SideToMove = fen.SideToMove;
            assert(EpTarget == Square::None); EpTarget   = fen.EPTarget;

            assert(DrawClock == 0); DrawClock  = fen.DrawClock;

            Rights = fen.CastlingRights;

            for (uint8 sq = 0; sq < Square::N; sq++)
                if (fen.Mailbox[sq] != ColoredPiece::None)
                    Insert(Square(sq), fen.Mailbox[sq]);

            GenerateCheckers();
        }

        constexpr inline ColoredPiece operator[](const Square sq) const {
            return Mailbox[static_cast<uint8>(sq)];
        }

        [[nodiscard]] constexpr std::string ToString() const {
            std::string board = "+---+---+---+---+---+---+---+---+\n";

            for (uint8 rank = 7; rank != 255; rank--) {
                board += "| ";

                for (uint8 file = 0; file < File::N; file++) {
                    board += Mailbox[rank * 8 + file].ToString() + " | ";
                }

                board += Rank(rank).ToString();
                board += "\n+---+---+---+---+---+---+---+---+\n";
            }

            board += "  a   b   c   d   e   f   g   h\n";
            return board;
        }

        constexpr inline void GenerateCheckers() {
            const auto friends = (*this)[ SideToMove];
            const auto enemies = (*this)[!SideToMove];
            const auto occupied = friends + enemies;

            //assert((*this)[Piece::King] != BitBoards::Empty);
            const auto king = ((*this)[Piece::King] & friends).LSB();
            //assert(king != Square::None);

            const auto p = (*this)[Piece::Pawn  ];
            const auto n = (*this)[Piece::Knight];
            const auto b = (*this)[Piece::Bishop];
            const auto r = (*this)[Piece::Rook  ];
            const auto q = (*this)[Piece::Queen ];

            const auto checkingP = p & MoveTable::Pawn(SideToMove, king);
            const auto checkingN = n & MoveTable::Knight(king);
            const auto checkingD = (b + q) & MoveTable::Bishop(king, occupied);
            const auto checkingL = (r + q) & MoveTable::Rook  (king, occupied);

            Checkers = (checkingP + checkingN + checkingD + checkingL) & enemies;
            CheckNum = Checkers.PopCount();
        }
    };

    constexpr inline std::ostream& operator<<(std::ostream& os, const Position& position) {
        os << position.ToString();
        return os;
    }
};

#endif //CHESS_POSITION
