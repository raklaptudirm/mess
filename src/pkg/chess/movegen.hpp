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
#include "position.hpp"

using namespace Chess;

namespace Chess {
    namespace Moves {
        template<Color STM, bool QUIET, bool NOISY>
        class Generator {
            const Position& position;
            const Castling::Info& castlingInfo;

            BitBoard friends;
            BitBoard enemies;
            BitBoard occupied;
            BitBoard blockers;

            BitBoard territory;

            Square king;

            BitBoard checkmask;

            BitBoard pinmaskL;
            BitBoard pinmaskD;

            MoveList moves;

            // serialize serializes the given targets BitBoard into an array
            // of moves from the given source square which are then appended
            // to the move-list.
            inline void serialize(Square source, BitBoard targets) {
                targets = targets & checkmask & territory;
                for (auto target : targets) moves += Move(source, target, Move::Flag::Normal);
            }

            template<Direction OFFSET, uint16 FLAG>
            inline void serialize(BitBoard targets) {
                targets = targets & checkmask & territory;
                for (auto target : targets) moves += Move(target >> -OFFSET, target, FLAG);
            }

            template<Direction OFFSET, bool CAPTURE>
            inline void serializePromotions(BitBoard targets) {
                targets = targets & checkmask & ~friends;
                for (auto target : targets) {
                    if (NOISY) moves += Move(target >> -OFFSET, target, Move::Flag::QPromotion);

                    if ((QUIET && !CAPTURE) || (NOISY && CAPTURE)) {
                        moves += Move(target >> -OFFSET, target, Move::Flag::NPromotion);
                        moves += Move(target >> -OFFSET, target, Move::Flag::BPromotion);
                        moves += Move(target >> -OFFSET, target, Move::Flag::RPromotion);
                    }
                }
            }

            inline void generateCheckMask() {
                switch (position.CheckNum) {
                    case 0: checkmask = BitBoards::Full;  break;
                    case 2: checkmask = BitBoards::Empty; break;
                    default:
                        const auto checkerBB = position.Checkers;
                        const auto checkerSq = checkerBB.LSB();
                        const auto checkerPc = position[checkerSq].Piece();

                        if (checkerPc == Piece::Pawn || checkerPc == Piece::Knight)
                             checkmask = checkerBB;
                        else checkmask = checkerBB + BitBoards::Between(king, checkerSq);
                }
            }

            inline void generatePinMasks() {
                const BitBoard b = enemies & position[Piece::Bishop];
                const BitBoard r = enemies & position[Piece::Rook  ];
                const BitBoard q = enemies & position[Piece::Queen ];

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

            inline void pawnMoves() {
                constexpr Direction UP = STM == Color::White ? Directions::North : Directions::South;
                constexpr Direction UE = UP + Directions::East, UW = UP + Directions::West;

                constexpr BitBoard DPRank = BitBoards::Rank(STM == Color::White ? Rank::Third  : Rank::Sixth);
                constexpr BitBoard PRRank = BitBoards::Rank(STM == Color::White ? Rank::Eighth : Rank::First);

                const BitBoard pawns = position[Piece::Pawn] & friends;

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

                    serializePromotions<UE, true>(attacksE & PRRank & enemies);
                    serializePromotions<UW, true>(attacksW & PRRank & enemies);

                    const Square epTarget = position.EpTarget;
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
                                            (position[Piece::Rook] + position[Piece::Queen]) & enemies;

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
                serializePromotions<UP, false>(singlePushes & PRRank);

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
                const BitBoard knights = (position[Piece::Knight] & friends) - (pinmaskL + pinmaskD);
                for (auto knight : knights) serialize(knight, MoveTable::Knight(knight));
            }

            // bishopMoves generates legal moves for bishop-like pieces, i.e. bishops and queens.
            inline void bishopMoves() {
                // Consider both bishops and queens. Pieces which are pinned
                // laterally can't make any diagonal moves, so Remove those.
                const BitBoard bishops = ((position[Piece::Bishop] + position[Piece::Queen]) & friends) - pinmaskL;

                // Pieces pinned diagonally can only make moves within
                // the pinned diagonal, so Remove all other targets.
                const BitBoard pinned = bishops & pinmaskD;
                for (auto bishop : pinned) serialize(bishop, MoveTable::Bishop(bishop, occupied) & pinmaskD);

                // Unpinned pieces can make any legal move.
                const BitBoard unpinned = bishops ^ pinned;
                for (auto bishop : unpinned) serialize(bishop, MoveTable::Bishop(bishop, occupied));
            }

            // rookMoves generates legal moves for rook-like pieces, i.e. rooks and queens.
            inline void rookMoves() {
                // Consider both rooks and queens. Pieces which are pinned
                // diagonally can't make any lateral moves, so Remove them.
                const BitBoard rooks = ((position[Piece::Rook] + position[Piece::Queen]) & friends) - pinmaskD;

                // Pieces pinned laterally can only make moves within
                // the pinned file/rank, so Remove all other targets.
                const BitBoard pinned = rooks & pinmaskL;
                for (auto rook : pinned) serialize(rook, MoveTable::Rook(rook, occupied) & pinmaskL);

                // Unpinned pieces can make any legal move.
                const BitBoard unpinned = rooks ^ pinned;
                for (auto rook : unpinned) serialize(rook, MoveTable::Rook(rook, occupied));
            }

            // kingMoves generates legal moves for the king.
            inline void kingMoves() {
                const BitBoard targets = MoveTable::King(king) & territory;

                for (auto target : targets) {
                    // Check if king move is legal.
                    if (!position.Attacked<!STM>(target, blockers))
                        moves += Move(king, target, Move::Flag::Normal);
                }
            }

            inline void castlingMoves() {
                if (!QUIET) return;

                #define GENERATE_CASTLING_MOVE(side)                                                               \
                {                                                                                                  \
                    const auto dimension = Castling::Dimension(STM, side);                                         \
                    if (                                                                                           \
                            /* Check for the necessary castling rights. */                                         \
                            position.Rights.Has(Castling::Rights(dimension)) &&                                    \
                            /* Check for blockers in the castling path. */                                         \
                            occupied.IsDisjoint(castlingInfo.BlockerMask(dimension)) &&                            \
                            /* Check for attackers in the king's path. */                                          \
                            !position.Attacked<!STM>(castlingInfo.AttackMask(dimension), blockers)                 \
                            ) {                                                                                    \
                        /* All castling requirements met: generate move */                                         \
                        moves += Move(king, castlingInfo.Rook(dimension), Move::Flag::FlagFrom(side));             \
                    }                                                                                              \
                } // GENERATE_CASTLING_MOVE

                GENERATE_CASTLING_MOVE(Castling::Side::H)
                GENERATE_CASTLING_MOVE(Castling::Side::A)

                #undef GENERATE_CASTLING_MOVE
            }

        public:

            explicit Generator(const Position& p, const Castling::Info& c) :
                    position(p), castlingInfo(c) {
                // Initialize various BitBoards.
                friends  = position[ STM];
                enemies  = position[!STM];
                occupied = friends + enemies;

                territory = BitBoards::Empty;
                if (QUIET) territory |= ~occupied;
                if (NOISY) territory |= enemies;

                const auto kingBB = position[Piece::King] & friends;

                blockers = occupied ^ kingBB;

                king = kingBB.LSB();

                generatePinMasks();
                generateCheckMask();
            }

            constexpr inline MoveList GenerateMoves() {
                moves.Clear();

                switch (position.CheckNum) {
                    case 0:
                        castlingMoves();
                    case 1:
                        rookMoves();
                        bishopMoves();
                        knightMoves();
                        pawnMoves();
                    default /* case 2 */:
                        kingMoves(); // Always generate king moves.
                }

                return moves;
            }
        };

        template<bool QUIET, bool NOISY>
        MoveList Generate(const Position& p, const Castling::Info& castlingInfo) {
            if (p.SideToMove == Color::White) {
                auto generator = Generator<Color::White, QUIET, NOISY>(p, castlingInfo);
                return generator.GenerateMoves();
            } else {
                auto generator = Generator<Color::Black, QUIET, NOISY>(p, castlingInfo);
                return generator.GenerateMoves();
            }
        }
    }
}

#endif