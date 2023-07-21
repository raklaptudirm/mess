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

#ifndef CHESS_MOVE_GENERATION
#define CHESS_MOVE_GENERATION

#include <iostream>

#include "move.hpp"
#include "board.hpp"
#include "moves.hpp"
#include "square.hpp"
#include "castling.hpp"
#include "bitboard.hpp"
#include "movelist.hpp"

using namespace Chess;

namespace Chess {
    namespace Moves {
        template<Color STM>
        class Generator {
            Board* board;
            Castling::Info* castling;

            BitBoard friends;
            BitBoard enemies;
            BitBoard occupied;
            BitBoard blockers;

            Square king;

            BitBoard checkmask;

            BitBoard pinmaskL;
            BitBoard pinmaskD;

            MoveList moves;

            // serialize serializes the given targets BitBoard into an array
            // of moves from the given source square which are then appended
            // to the move-list.
            inline void serialize(Square source, BitBoard targets) {
                targets = targets & checkmask & ~friends;
                for (auto target : targets) moves += Move(source, target, Move::Flag::Normal);
            }

            template<Direction OFFSET, uint16 FLAG>
            inline void serialize(BitBoard targets) {
                targets = targets & checkmask & ~friends;
                for (auto target : targets) moves += Move(target >> -OFFSET, target, FLAG);
            }

            template<Direction OFFSET, bool QUIET, bool NOISY>
            inline void serializePromotions(BitBoard targets) {
                targets = targets & checkmask & ~friends;
                for (auto target : targets) {
                    if (NOISY) moves += Move(target >> -OFFSET, target, Move::Flag::QPromotion);

                    if (QUIET) {
                        moves += Move(target >> -OFFSET, target, Move::Flag::NPromotion);
                        moves += Move(target >> -OFFSET, target, Move::Flag::BPromotion);
                        moves += Move(target >> -OFFSET, target, Move::Flag::RPromotion);
                    }
                }
            }

            inline void generateCheckMask() {
                switch (board->CheckNum()) {
                    case 0: checkmask = BitBoards::Full;  break;
                    case 2: checkmask = BitBoards::Empty; break;
                    default:
                        const auto checkerBB = board->Checkers();
                        const auto checkerSq = checkerBB.LSB();
                        const auto checkerPc = (*board)[checkerSq].Piece();

                        if (checkerPc == Piece::Pawn || checkerPc == Piece::Knight)
                             checkmask = checkerBB;
                        else checkmask = checkerBB + BitBoards::Between(king, checkerSq);
                }
            }

            inline void generatePinMasks() {
                const BitBoard b = enemies & (*board)[Piece::Bishop];
                const BitBoard r = enemies & (*board)[Piece::Rook  ];
                const BitBoard q = enemies & (*board)[Piece::Queen ];

                // Fetch the possibly pinning Bishops, Rooks, and Queens.
                const BitBoard pinningL = (r | q) & MoveTable::Rook  (king, enemies);
                const BitBoard pinningD = (b | q) & MoveTable::Bishop(king, enemies);

                // Empty the Pinmasks.
                pinmaskL = BitBoards::Empty;
                pinmaskD = BitBoards::Empty;

                /******************************
                 * Lateral Pinmask Generation *
                 ******************************/
                for (const auto& piece : pinningL) {
                    const BitBoard possiblePin = BitBoards::Between(king, piece);

                    if ((friends & possiblePin).Singular())
                        pinmaskL = pinmaskL | possiblePin | BitBoard(piece);
                }

                /*******************************
                 * Diagonal Pinmask Generation *
                 *******************************/
                for (const auto& piece : pinningD) {
                    const BitBoard possiblePin = BitBoards::Between(king, piece);

                    if ((friends & possiblePin).Singular())
                        pinmaskD = pinmaskD | possiblePin | BitBoard(piece);
                }
            }

            template<bool QUIET, bool NOISY>
            inline void pawnMoves() {
                constexpr Direction UP = STM == Color::White ? Directions::North : Directions::South;
                constexpr Direction UE = UP + Directions::East, UW = UP + Directions::West;

                constexpr BitBoard DPRank = BitBoards::Rank(STM == Color::White ? Rank::Third  : Rank::Sixth);
                constexpr BitBoard PRRank = BitBoards::Rank(STM == Color::White ? Rank::Eighth : Rank::First);

                const BitBoard pawns = (*board)[Piece::Pawn] & friends;

                /****************************
                 * Pawn Captures Generation *
                 ****************************/
                if (NOISY) {
                    const BitBoard attackers = pawns - pinmaskL;

                    const BitBoard   pinnedAttackers = attackers & pinmaskD;
                    const BitBoard unpinnedAttackers = attackers ^ pinnedAttackers;

                    const BitBoard   pinnedAttacksE = (  pinnedAttackers >> UE);
                    const BitBoard   pinnedAttacksW = (  pinnedAttackers >> UW);
                    const BitBoard unpinnedAttacksE = (unpinnedAttackers >> UE);
                    const BitBoard unpinnedAttacksW = (unpinnedAttackers >> UW);

                    const BitBoard attacksE = (pinnedAttacksE & pinmaskD) | unpinnedAttacksE;
                    const BitBoard attacksW = (pinnedAttacksW & pinmaskD) | unpinnedAttacksW;

                    serialize<UE, Move::Flag::Normal>((attacksE - PRRank) & enemies);
                    serialize<UW, Move::Flag::Normal>((attacksW - PRRank) & enemies);

                    serializePromotions<UE, QUIET, NOISY>(attacksE & PRRank & enemies);
                    serializePromotions<UW, QUIET, NOISY>(attacksW & PRRank & enemies);

                    const Square epTarget = board->EPTarget();
                    if (epTarget != Square::None) {
                        const BitBoard target = BitBoard(epTarget);
                        const BitBoard passanters = MoveTable::Pawn<!STM>(epTarget) & attackers;

                        switch (passanters.PopCount()) {
                            case 1: {
                                if ((target + (target >> -UP)).IsDisjoint(checkmask)) {
                                    break;
                                }

                                const auto captured = epTarget >> -UP;
                                if (king.Rank() == captured.Rank()) {
                                    const BitBoard pinners =
                                        ((*board)[Piece::Rook] + (*board)[Piece::Queen]) & enemies;

                                    const BitBoard vanishers = passanters + BitBoard(captured);

                                    if (!MoveTable::Rook(king, occupied ^ vanishers).IsDisjoint(pinners))
                                        break;
                                }

                                if (pinmaskD.IsDisjoint(passanters) || !pinmaskD.IsDisjoint(target))
                                    moves += Move(passanters.LSB(), epTarget, Move::Flag::EnPassant);

                                break;
                            }

                            case 2: {
                                for (const auto passanter : passanters)
                                    if (!pinmaskD[passanter] || !pinmaskD.IsDisjoint(target))
                                        moves += Move(passanter, epTarget, Move::Flag::EnPassant);
                            }
                        }
                    }
                }

                const BitBoard pushers = pawns - pinmaskD;

                const BitBoard   pinnedPushers = pushers & pinmaskL;
                const BitBoard unpinnedPushers = pushers ^ pinnedPushers;

                const BitBoard   pinnedSinglePush = ((  pinnedPushers >> UP) - occupied);
                const BitBoard unpinnedSinglePush = ((unpinnedPushers >> UP) - occupied);

                const BitBoard singlePushes = ((pinnedSinglePush & pinmaskL) + unpinnedSinglePush);

                /*****************************
                 * Promotion Push Generation *
                 *****************************/
                serializePromotions<UP, QUIET, NOISY>(singlePushes & PRRank);

                /****************************************
                 * Normal Single/Double Push Generation *
                 ****************************************/
                if (QUIET) {
                    const BitBoard doublePushes = ((singlePushes & DPRank) >> UP) - occupied;

                    serialize<   UP, Move::Flag::Normal    >(singlePushes - PRRank); // Single Pushes.
                    serialize<UP+UP, Move::Flag::DoublePush>(doublePushes);          // Double Pushes.
                }
            }

            // knightMoves generates legal moves for knights.
            inline void knightMoves() {
                // Knights which are pinned either laterally or diagonally can't move.
                const BitBoard knights = ((*board)[Piece::Knight] & friends) - (pinmaskL + pinmaskD);
                for (auto knight : knights) serialize(knight, MoveTable::Knight(knight));
            }

            // bishopMoves generates legal moves for bishop-like pieces, i.e. bishops and queens.
            inline void bishopMoves() {
                // Consider both bishops and queens. Pieces which are pinned
                // laterally can't make any diagonal moves, so remove those.
                const BitBoard bishops = (((*board)[Piece::Bishop] + (*board)[Piece::Queen]) & friends) - pinmaskL;

                // Pieces pinned diagonally can only make moves within
                // the pinned diagonal, so remove all other targets.
                const BitBoard pinned = bishops & pinmaskD;
                for (auto bishop : pinned) serialize(bishop, MoveTable::Bishop(bishop, occupied) & pinmaskD);

                // Unpinned pieces can make any legal move.
                const BitBoard unpinned = bishops ^ pinned;
                for (auto bishop : unpinned) serialize(bishop, MoveTable::Bishop(bishop, occupied));
            }

            // rookMoves generates legal moves for rook-like pieces, i.e. rooks and queens.
            inline void rookMoves() {
                // Consider both rooks and queens. Pieces which are pinned
                // diagonally can't make any lateral moves, so remove them.
                const BitBoard rooks = (((*board)[Piece::Rook] + (*board)[Piece::Queen]) & friends) - pinmaskD;

                // Pieces pinned laterally can only make moves within
                // the pinned file/rank, so remove all other targets.
                const BitBoard pinned = rooks & pinmaskL;
                for (auto rook : pinned) serialize(rook, MoveTable::Rook(rook, occupied) & pinmaskL);

                // Unpinned pieces can make any legal move.
                const BitBoard unpinned = rooks ^ pinned;
                for (auto rook : unpinned) serialize(rook, MoveTable::Rook(rook, occupied));
            }

            // kingMoves generates legal moves for the king.
            inline void kingMoves() {
                const BitBoard targets = MoveTable::King(king) & ~friends;

                for (auto target : targets) {
                    // Check if king move is legal.
                    if (!(*board).Attacked<!STM>(target, blockers))
                        moves += Move(king, target, Move::Flag::Normal);
                }
            }

            inline void castlingMoves() {
                const auto pathH = castling->PathH<STM>();
                if (
                    castling->Rights.Has(Castling::H<STM>()) &&
                    occupied.IsDisjoint(pathH) &&
                    !board->Attacked<!STM>(pathH, blockers)
                ) {
                    moves += Move(king, castling->RookH<STM>(), Move::Flag::CastleHSide);
                }

                const auto pathA = castling->PathA<STM>();
                if (
                    castling->Rights.Has(Castling::A<STM>()) &&
                    occupied.IsDisjoint(pathA + (pathA >> Directions::West)) &&
                    !board->Attacked<!STM>(pathA, blockers)
                ) {
                    moves += Move(king, castling->RookA<STM>(), Move::Flag::CastleASide);
                }
            }

        public:

            Generator(Board* b, Castling::Info* c) :
                board(b), castling(c) {
                UpdateContext(b, c);
            }

            void UpdateContext(Board* b, Castling::Info* c) {
                board = b;
                castling = c;

                // Initialize various BitBoards.
                friends  = (*board)[ STM];
                enemies  = (*board)[!STM];
                occupied = friends + enemies;

                const auto kingBB = (*board)[Piece::King] & friends;

                blockers = occupied ^ kingBB;

                king = kingBB.LSB();

                generatePinMasks();
                generateCheckMask();
            }

            template<bool QUIET, bool NOISY>
            constexpr inline MoveList GenerateMoves() {
                moves.Clear();

                switch (board->CheckNum()) {
                    case 0:
                        castlingMoves();
                    case 1:
                        rookMoves();
                        bishopMoves();
                        knightMoves();
                        pawnMoves<QUIET, NOISY>();
                    default /* case 2 */:
                        kingMoves(); // Always generate king moves.
                }

                return moves;
            }
        };

        template<bool QUIET, bool NOISY>
        MoveList Generate(Board *b) {
            if (b->SideToMove() == Color::White) {
                static auto generator = Generator<Color::White>(b, b->Castling());

                generator.UpdateContext(b, b->Castling());
                return generator.GenerateMoves<QUIET, NOISY>();
            } else {
                static auto generator = Generator<Color::Black>(b, b->Castling());

                generator.UpdateContext(b, b->Castling());
                return generator.GenerateMoves<QUIET, NOISY>();
            }
        }
    }
}

#endif