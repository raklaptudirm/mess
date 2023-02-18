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

package board

import (
	"strconv"
	"strings"

	"laptudirm.com/x/mess/pkg/board/move/castling"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
	"laptudirm.com/x/mess/pkg/board/zobrist"
	"laptudirm.com/x/mess/pkg/formats/fen"
)

var StartFEN = fen.FromString("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

// NewBoard creates an instance of a *Board from a given fen string.
// https://www.chessprogramming.org/Forsyth-Edwards_Notation
func (board *Board) UpdateWithFEN(fen fen.String) {
	// empty board
	for s := square.A8; s <= square.H1; s++ {
		board.ClearSquare(s)
	}

	// reset some stuff
	board.Plys = 0
	board.Hash = 0

	// side to move
	board.SideToMove = piece.NewColor(fen[1])
	if board.SideToMove == piece.Black {
		board.Hash ^= zobrist.SideToMove
	}

	// generate position
	ranks := strings.Split(fen[0], "/")
	for rankId, rankData := range ranks {
		fileId := square.FileA
		for _, id := range rankData {
			s := square.New(fileId, square.Rank(rankId))

			if id >= '1' && id <= '8' {
				skip := square.File(id - 48) // ascii value to number
				fileId += skip               // skip over squares
				continue
			}

			// piece string to piece
			p := piece.NewFromString(string(id))

			if p.Type() != piece.NoType {
				board.FillSquare(s, p)
			}

			fileId++
		}
	}

	// castling rights
	board.CastlingRights = castling.NewRights(fen[2])
	board.Hash ^= zobrist.Castling[board.CastlingRights]

	// en-passant target square
	board.EnPassantTarget = square.NewFromString(fen[3])
	if board.EnPassantTarget != square.None {
		board.Hash ^= zobrist.EnPassant[board.EnPassantTarget.File()]
	}

	// move counters
	board.DrawClock, _ = strconv.Atoi(fen[4])
	board.FullMoves, _ = strconv.Atoi(fen[5])
}

// FEN returns the fen string of the current Board position.
func (b *Board) FEN() fen.String {
	var fenString fen.String
	fenString[0] = b.Position.FEN()
	fenString[1] = b.SideToMove.String()
	fenString[2] = b.CastlingRights.String()
	fenString[3] = b.EnPassantTarget.String()
	fenString[4] = strconv.Itoa(b.DrawClock)
	fenString[5] = strconv.Itoa(b.FullMoves)
	return fenString
}
