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

#include "move.hpp"
#include "board.hpp"
#include "moves.hpp"
#include "square.hpp"
#include "castling.hpp"
#include "bitboard.hpp"
#include "movelist.hpp"
#include "position.hpp"

using namespace Chess;

namespace Chess::Moves {
    template<Color STM, bool QUIET, bool NOISY>
    class Generator {
        // Position information, along with the relevant
        // CastlingInfo used for generating castling moves.
        const Position& position;
        const Castling::Info& castlingInfo;

        // BitBoards for various board features.
        BitBoard friends;
        BitBoard enemies;
        BitBoard occupied;

        // Blockers is the occupied BitBoard without
        // the side to move's king in the set. Useful
        // when calculating safe squares for the king.
        BitBoard blockers;

        // BitBoard which represents the squares to which
        // moves can be made. Its value is dictated by the
        // QUIET and NOISY variables, which allows/disallows
        // movement to empty squares and enemy occupied
        // squares respectively. This is not applicable to
        // promotions as there other metrics may be used
        // for determining if a move is quiet or noisy.
        BitBoard territory;

        // BitBoard which is kind of like territory, except
        // its value is dictated by if and how the king is
        // being checked. If the king is under check, it
        // marks the squares moving to which will block the
        // check, otherwise it contains all the squares. It
        // is stored separately from territory cause checkmask
        // doesn't restrain the king's movement, while territory
        // does restrain it.
        BitBoard checkmask;

        // Square of the side to move's king.
        Square king;

        // Lateral and Diagonal pinmasks are BitBoards which
        // store friendly pieces which are pinned laterally or
        // diagonally, along with the rays along which they are
        // pinned. This allows us to restrain the pinned piece's
        // move along the ray as otherwise the king would be in
        // check and open to be captured by the enemy.
        BitBoard pinmaskL;
        BitBoard pinmaskD;

        // The internal movelist which stores all the moves.
        MoveList moves;

        // serialize serializes the given targets BitBoard into an array
        // of moves from the given source square which are then appended
        // to the move-list.
        inline void serialize(Square source, BitBoard targets) {
            targets = targets & checkmask & territory;
            for (auto target : targets) moves += Move(source, target, Move::Flag::Normal);
        }

        // Overload of serialize which infers the source square from the
        // target square and the target-source offset. It also accepts a
        // move flag which is packed into the final move.
        template<Direction OFFSET, uint16 FLAG>
        inline void serialize(BitBoard targets) {
            targets = targets & checkmask & territory;
            for (auto target : targets) moves += Move(target >> -OFFSET, target, FLAG);
        }

        // serializePromotions is similar to the offset overload of serialize
        // as it also infers the source from the target and the target-source
        // offset. It additionally generates all the possible promotion types
        // according to the provided move generation type.
        template<Direction OFFSET, bool CAPTURE>
        inline void serializePromotions(BitBoard targets) {
            // Unlike other serialization methods, the target BitBoard is not
            // masked with territory since queen promotions are noisy moves
            // which may move to empty squares. Therefore, the territory
            // logic is implemented inside the target loop.
            targets = targets & checkmask & ~friends;
            for (auto target : targets) {
                // Queen promotions are noisy moves, so generate them whenever
                // we can generate noisy moves according to the generation type.
                if (NOISY) moves += Move(target >> -OFFSET, target, Move::Flag::QPromotion);

                // Other types of promotions are quiet moves by default, so
                // their noisy-ness is determined like that of any other move:
                // whether they are a capture or a non-capture.
                if ((QUIET && !CAPTURE) || (NOISY && CAPTURE)) {
                    moves += Move(target >> -OFFSET, target, Move::Flag::NPromotion);
                    moves += Move(target >> -OFFSET, target, Move::Flag::BPromotion);
                    moves += Move(target >> -OFFSET, target, Move::Flag::RPromotion);
                }
            }
        }

        // generateCheckMask generates the checkmask for the current position.
        // Look at the documentation for the checkmask variable for more info.
        inline void generateCheckMask() {
            switch (position.CheckNum) {
                // King is not under any checks, all moves are possible.
                case 0: checkmask = BitBoards::Full;  break;

                // King is under double check, no moves are possible for non-king pieces.
                case 2: checkmask = BitBoards::Empty; break;

                // King is under a singular check. Determine the type of check and set
                // the value of the checkmask accordingly.
                default:
                    const auto checkerSq = position.Checkers.LSB();
                    const auto checkerPc = position[checkerSq].Piece();

                    if (checkerPc == Piece::Pawn || checkerPc == Piece::Knight)
                        // Pawn/Knight checks cannot be blocked. Only possible moves
                        // by non-king pieces is capturing the checking piece.
                        checkmask = position.Checkers;
                    else // Sliding piece moves can be blocked, so include the between
                        // squares in the checkmask along with the checking piece.
                        checkmask = BitBoards::Between2(king, checkerSq);
            }
        }

        [[nodiscard]] inline BitBoard generatePinMask(const BitBoard pinning) const {
            auto pinmask = BitBoards::Empty;
            for (const auto& piece : pinning) {
                // Get the possibly pinning ray (can have friendly pieces).
                const BitBoard possiblePin = BitBoards::Between2(king, piece);

                // If the number of friendly pieces in the pinning ray is exactly
                // one, then that piece is being pinned along that ray.
                if ((friends & possiblePin).Singular())
                    pinmask |= possiblePin;
            }

            return pinmask;
        }

        // generatePinMasks generates the lateral and diagonal pinmasks.
        // Look at the documentation for pinmaskL/pinmaskD for more info.
        inline void generatePinMasks() {
            // Get enemy sliding pieces, which can pin pieces.
            const BitBoard b = enemies & position[Piece::Bishop];
            const BitBoard r = enemies & position[Piece::Rook  ];
            const BitBoard q = enemies & position[Piece::Queen ];

            // Fetch the possibly pinning Bishops, Rooks, and Queens: the ones whose attacks line
            // up with the position of the side to move's king, and generate the pinmasks.
            pinmaskL = generatePinMask((r | q) & MoveTable::Rook  (king, enemies));
            pinmaskD = generatePinMask((b | q) & MoveTable::Bishop(king, enemies));
        }

        // pawnMoves generates all the different types of pawn moves that are legal
        // in this position and are in accordance with the move generation type.
        inline void pawnMoves() {
            // Some useful direction constants.
            constexpr Direction UP = STM == Color::White ? Directions::North : Directions::South;
            constexpr Direction UE = UP + Directions::East, UW = UP + Directions::West;

            // Some useful rank BitBoard constants including the double push and promotion ranks.
            constexpr BitBoard DPRank = BitBoards::Rank(STM == Color::White ? Rank::Third  : Rank::Sixth);
            constexpr BitBoard PRRank = BitBoards::Rank(STM == Color::White ? Rank::Eighth : Rank::First);

            // BitBoard containing all friendly pawns whose moves we are generating.
            const BitBoard pawns = position[Piece::Pawn] & friends;

            /* **************************
             * Pawn Captures Generation *
             ************************** */
            if (NOISY) { // Only generate captures if noisy moves are allowed.

                // Captures are diagonal moves so pawns pinned laterally can't capture.
                const BitBoard attackers = pawns - pinmaskL;

                // Separate the pawns into groups depending on whether they are pinned
                // diagonally or not. A pawn which is pinned diagonally can only move
                // in the pinned direction.
                const BitBoard   pinnedAttackers = attackers & pinmaskD;
                const BitBoard unpinnedAttackers = attackers ^ pinnedAttackers;

                // Shift the pawns into the squares they are attacking.
                const BitBoard   pinnedAttacksE = (  pinnedAttackers >> UE);
                const BitBoard   pinnedAttacksW = (  pinnedAttackers >> UW);
                const BitBoard unpinnedAttacksE = (unpinnedAttackers >> UE);
                const BitBoard unpinnedAttacksW = (unpinnedAttackers >> UW);

                // Concatenate the attacks of the pinned and the unpinned pawns into
                // singular variables in each direction. Notice we do an intersection
                // of the pinned attacks and the pinmask to remove illegal moves.
                const BitBoard attacksE = (pinnedAttacksE & pinmaskD) | unpinnedAttacksE;
                const BitBoard attacksW = (pinnedAttacksW & pinmaskD) | unpinnedAttacksW;

                // Serialize the non-promotion attacks which actually capture an enemy.
                serialize<UE, Move::Flag::Normal>((attacksE - PRRank) & enemies);
                serialize<UW, Move::Flag::Normal>((attacksW - PRRank) & enemies);

                // Serialize the promotion captures.
                serializePromotions<UE, true>(attacksE & PRRank & enemies);
                serializePromotions<UW, true>(attacksW & PRRank & enemies);

                /* ***********************
                 * En Passant Generation *
                 *********************** */
                if (position.EpTarget != Square::None) {
                    const auto target = position.EpTarget;
                    const auto targetBB = BitBoard(target);

                    // BitBoard containing friendly pawns which attack the target.
                    const auto passanters = MoveTable::Pawn<!STM>(target) & attackers;

                    switch (passanters.PopCount()) {
                        // Only one passanter: possible king double pin.
                        case 1: {
                            if ((targetBB + (targetBB >> -UP)).IsDisjoint(checkmask)) {
                                break;
                            }

                            const auto captured = target >> -UP;
                            if (king.Rank() == captured.Rank()) {
                                const BitBoard pinners =
                                        (position[Piece::Rook] + position[Piece::Queen]) & enemies;

                                const BitBoard vanishers = passanters + BitBoard(captured);

                                if (!MoveTable::Rook(king, occupied ^ vanishers).IsDisjoint(pinners))
                                    break;
                            }

                            if (pinmaskD.IsDisjoint(passanters) || !pinmaskD.IsDisjoint(targetBB))
                                moves += Move(passanters.LSB(), target, Move::Flag::EnPassant);

                            break;
                        }

                        // Two passanters, king double pin is impossible so simply iterate
                        // over the passanters and generate the legal en-passant moves.
                        case 2: {
                            for (const auto passanter : passanters)
                                if (!pinmaskD[passanter] || !pinmaskD.IsDisjoint(targetBB))
                                    moves += Move(passanter, target, Move::Flag::EnPassant);
                        }
                    }
                }
            }

            /* ************************************
             * Pawn Single/Double Push Generation *
             ************************************ */
            if (QUIET) {
                // Pushes are lateral moves so diagonally pinned pawns can't push.
                const BitBoard pushers = pawns - pinmaskD;

                // Separate the pawns into groups depending on whether they are pinned
                // laterally or not. A pawn which is pinned laterally can only move
                // in the pinned direction.
                const BitBoard   pinnedPushers = pushers & pinmaskL;
                const BitBoard unpinnedPushers = pushers ^ pinnedPushers;

                // Shift the pawns up into their target squares, removing the ones which
                // collide with other pieces to get all the single pushes.
                const BitBoard   pinnedSinglePush = ((  pinnedPushers >> UP) - occupied);
                const BitBoard unpinnedSinglePush = ((unpinnedPushers >> UP) - occupied);

                // Combine the pinned and unpinned single pushes into a single BitBoard.
                const BitBoard singlePushes = ((pinnedSinglePush & pinmaskL) + unpinnedSinglePush);

                // Push the single pushed from the double push rank upwards and remove
                // the ones which collide with other pieces to get the double pushes.
                const BitBoard doublePushes = ((singlePushes & DPRank) >> UP) - occupied;

                // Serialize the single and double pushes. Remove the promotion rank
                // from the serialization of the single pushes as they are handled
                // separately so that all the promotions are properly generated.
                serialize<   UP, Move::Flag::Normal    >(singlePushes - PRRank); // Single Pushes.
                serialize<UP+UP, Move::Flag::DoublePush>(doublePushes);          // Double Pushes.

                // Serialize the promotions by extracting the pushes in the promotion rank.
                serializePromotions<UP, false>(singlePushes & PRRank);
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
            // laterally can't make any diagonal moves, so remove those.
            const BitBoard bishops = ((position[Piece::Bishop] + position[Piece::Queen]) & friends) - pinmaskL;

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
            const BitBoard rooks = ((position[Piece::Rook] + position[Piece::Queen]) & friends) - pinmaskD;

            // Pieces pinned laterally can only make moves within
            // the pinned file/rank, so remove all other targets.
            const BitBoard pinned = rooks & pinmaskL;
            for (auto rook : pinned) serialize(rook, MoveTable::Rook(rook, occupied) & pinmaskL);

            // Unpinned pieces can make any legal move.
            const BitBoard unpinned = rooks ^ pinned;
            for (auto rook : unpinned) serialize(rook, MoveTable::Rook(rook, occupied));
        }

        // kingMoves generates legal moves for the king, excluding castling.
        inline void kingMoves() {
            const BitBoard targets = MoveTable::King(king) & territory;

            for (auto target : targets) {
                // Check if king move is legal.
                if (!position.Attacked<!STM>(target, blockers))
                    moves += Move(king, target, Move::Flag::Normal);
            }
        }

        // castlingMove tries to generate a castling move for the given side.
        inline void castlingMove(Castling::Side side) {
            const auto dimension = Castling::Dimension(STM, side);
            if (    // Check if castling requirements are met:
                    position.Rights.Has(Castling::Rights(dimension)) && // Check for the necessary castling rights.
                    occupied.IsDisjoint(castlingInfo.BlockerMask(dimension)) && // Check for blockers in the castling path.
                    !position.Attacked<!STM>(castlingInfo.AttackMask(dimension), blockers) // Check for attackers in the king's path.
                ) {
                moves += Move(king, castlingInfo.Rook(dimension), Move::Flag::FlagFrom(side));
            }
        }

        // castlingMoves generates all legal castling moves.
        inline void castlingMoves() {
            // Generate castling only if quiet moves are allowed.
            if (!QUIET) return;

            // Try to generate castling move for both sides.
            castlingMove(Castling::Side::H);
            castlingMove(Castling::Side::A);
        }

    public:

        // NOLINTNEXTLINE(cppcoreguidelines-pro-type-member-init)
        Generator(const Position& p, const Castling::Info& c) :
                position(p), castlingInfo(c) {

            // Initialize various BitBoards.
            friends  = position[ STM];
            enemies  = position[!STM];
            occupied = friends + enemies;

            // Initialize the territory BitBoard.
            territory = BitBoards::Empty;
            if (QUIET) territory |= ~occupied; // QUIET => Can move to empty squares.
            if (NOISY) territory |= enemies;   // NOISY => Can move to enemy squares.

            const auto kingBB = position[Piece::King] & friends;

            // Generate blockers bitboard (occupied - kingBB).
            blockers = occupied ^ kingBB;

            // Store the side to move's king's square.
            king = kingBB.LSB();

            generatePinMasks();
            generateCheckMask();
        }

        constexpr inline MoveList GenerateMoves() {
            // Clear the move list for move generation.
            moves.Clear();

            // Due to the fallthrough nature of switch statements, the move generation
            // functions below get called when the matched case is the same as or higher
            // that the case that contains the function. In this particular case, this
            // implies that the functions get called when the number of checks is less
            // than or equal to their threshold.
            switch (position.CheckNum) {
                case 0:
                    // Castling is only possible if
                    // the king is not in check.
                    castlingMoves();
                case 1:
                    // Non-king moves are only possible
                    // if the king is not in double check.
                    rookMoves();
                    bishopMoves();
                    knightMoves();
                    pawnMoves();
                default /* case 2 */:
                    // King moves are always possible.
                    kingMoves();
            }

            // Return a copy of the generated movelist.
            return moves;
        }
    };

    // Generate generates all the possible legal moves on the given Position and
    // with the given CastlingInfo which match the provided move generation type.
    template<bool QUIET, bool NOISY>
    MoveList Generate(const Position& p, const Castling::Info& castlingInfo) {
        // Switch template arguments according to side to move.
        if (p.SideToMove == Color::White) {
            auto generator = Generator<Color::White, QUIET, NOISY>(p, castlingInfo);
            return generator.GenerateMoves();
        } else {
            auto generator = Generator<Color::Black, QUIET, NOISY>(p, castlingInfo);
            return generator.GenerateMoves();
        }
    }
}

#endif // CHESS_MOVE_GENERATION