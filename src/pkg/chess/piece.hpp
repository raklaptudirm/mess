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

#ifndef CHESS_PIECE
#define CHESS_PIECE

#include <string>
#include <cassert>

#include "../util/types.hpp"

#include "piece.hpp"
#include "color.hpp"

namespace Chess {
    struct Piece {

        /********************************
         * Internal Enum Representation *
         ********************************/

        static const int N = 6;

        enum internal_type : uint8 {
            Pawn, Knight, Bishop, Rook, Queen, King, None
        };

        internal_type internal = None;

        /***************************
         * Constructor Definitions *
         ***************************/

        constexpr inline Piece(internal_type piece) {
            internal = piece;
        }

        constexpr explicit inline Piece(uint8 piece) {
            internal = static_cast<internal_type>(piece);
        }

        [[nodiscard]] constexpr inline std::string ToString() const {
            const std::string runes = "pnbrqk-";
            return runes.substr(static_cast<uint8>(internal), 1);
        }

        /************************
         * Conversion Functions *
         ************************/

        constexpr inline explicit operator uint8() const {
            return static_cast<uint8>(internal);
        }

        constexpr inline bool operator ==(const Piece&) const = default;
    };

    struct ColoredPiece {

        /********************************
         * Internal Enum Representation *
         ********************************/

        static const int N = 12;

        enum internal_type : uint8 {
            // White Pieces.
            WhitePawn, WhiteKnight, WhiteBishop,
            WhiteRook, WhiteQueen, WhiteKing,

            // Black Pieces.
            BlackPawn, BlackKnight, BlackBishop,
            BlackRook, BlackQueen, BlackKing,

            None,
        };

        internal_type internal = None;

        /***************************
         * Constructor Definitions *
         ***************************/

        constexpr inline ColoredPiece() {
            internal = ColoredPiece::None;
        }

        constexpr inline ColoredPiece(internal_type piece) {
            internal = piece;
        }

        constexpr inline ColoredPiece(Piece piece, Color color) {
            internal = static_cast<internal_type>(static_cast<uint8>(color)*Piece::N + static_cast<uint8>(piece));
        }

        constexpr inline explicit ColoredPiece(const std::string& piece) {
            assert(piece.length() == 1);

            if (piece == "P") internal = WhitePawn;
            else if (piece == "N") internal = WhiteKnight;
            else if (piece == "B") internal = WhiteBishop;
            else if (piece == "R") internal = WhiteRook;
            else if (piece == "Q") internal = WhiteQueen;
            else if (piece == "K") internal = WhiteKing;

            else if (piece == "p") internal = BlackPawn;
            else if (piece == "n") internal = BlackKnight;
            else if (piece == "b") internal = BlackBishop;
            else if (piece == "r") internal = BlackRook;
            else if (piece == "q") internal = BlackQueen;
            else if (piece == "k") internal = BlackKing;
        }

        [[nodiscard]] constexpr inline std::string ToString() const {
            const std::string runes = "PNBRQKpnbrqk-";
            return runes.substr(static_cast<uint8>(internal), 1);
        }

        /*****************************
         * Property Getter Functions *
         *****************************/

        [[nodiscard]] constexpr inline Piece Piece() const {
            if (internal == None) return Piece::None;
            return Chess::Piece(static_cast<uint8>(internal) % Piece::N);
        }

        [[nodiscard]] constexpr inline Color Color() const {
            return Chess::Color(static_cast<uint8>(internal) / Piece::N);
        }

        /************************
         * Conversion Functions *
         ************************/

        constexpr inline explicit operator uint8() const {
            return static_cast<uint8>(internal);
        }

        constexpr inline bool operator ==(const ColoredPiece&) const = default;
    };

    constexpr inline ColoredPiece operator +(Piece piece, Color color) {
        return ColoredPiece(piece, color);
    }
}

#endif