// Copyright © 2022 Rak Laptudirm <raklaptudirm@gmail.com>
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

package board

import (
	"fmt"
	"strconv"
	"strings"

	"laptudirm.com/x/mess/pkg/move"
	"laptudirm.com/x/mess/pkg/piece"
	"laptudirm.com/x/mess/pkg/square"
	"laptudirm.com/x/mess/pkg/zobrist"
)

// New creates an instance of a *Board from a given fen string.
// https://www.chessprogramming.org/Forsyth-Edwards_Notation
func New(fen string) *Board {
	var board Board

	parts := strings.Split(fen, " ")

	// side to move
	board.sideToMove = piece.NewColor(parts[1])
	if board.sideToMove == piece.BlackColor {
		board.hash ^= zobrist.SideToMove
	}

	// generate position
	ranks := strings.Split(parts[0], "/")
	for rankId, rankData := range ranks {
		fileId := square.FileA
		for _, id := range rankData {
			s := square.From(fileId, square.Rank(rankId))

			if id >= '1' && id <= '8' {
				skip := square.File(id - 48) // ascii value to number
				fileId += skip               // skip over squares
				continue
			}

			// piece string to piece
			p := piece.New(string(id))

			if p.Type() != piece.Empty {
				// update hash
				board.hash ^= zobrist.PieceSquare[p][s]

				// update board
				board.position[s] = p     // 8x8
				board.bitboards[p].Set(s) // bitboard

				// update friend and enemy bitboards
				if p.Color() == board.sideToMove {
					board.friends.Set(s)
				} else {
					board.enemies.Set(s)
				}
			}

			fileId++
		}
	}

	// castling rights
	board.castlingRights = move.CastlingRightsFrom(parts[2])
	board.hash ^= zobrist.Castling[board.castlingRights]

	// en-passant target square
	board.enPassantTarget = square.New(parts[3])
	if board.enPassantTarget != square.None {
		board.hash ^= zobrist.EnPassant[board.enPassantTarget.File()]
	}

	// move counters
	board.halfMoves, _ = strconv.Atoi(parts[4])
	board.fullMoves, _ = strconv.Atoi(parts[5])

	return &board
}

// FEN returns the fen string of the current Board position.
func (b *Board) FEN() string {
	// castling rights
	var castling string
	if castling = b.castlingRights.String(); castling != "" {
		castling += " "
	}

	// <position> <side to move> <castling rights> <en passant target> <half move count> <full move count>
	return fmt.Sprintf("%s %s %s%s %d %d", b.position.FEN(), b.sideToMove, castling, b.enPassantTarget, b.halfMoves, b.fullMoves)
}
