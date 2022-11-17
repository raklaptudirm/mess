// Copyright © 2022 Rak Laptudirm <rak@laptudirm.com>
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

package chess

import (
	"strconv"
	"strings"

	"laptudirm.com/x/mess/pkg/chess/move/castling"
	"laptudirm.com/x/mess/pkg/chess/piece"
	"laptudirm.com/x/mess/pkg/chess/square"
	"laptudirm.com/x/mess/pkg/zobrist"
)

// NewBoard creates an instance of a *Board from a given fen string.
// https://www.chessprogramming.org/Forsyth-Edwards_Notation
func NewBoard(fen string) *Board {
	var board Board

	parts := strings.Split(fen, " ")

	// side to move
	board.SideToMove = piece.NewColor(parts[1])
	if board.SideToMove == piece.Black {
		board.Hash ^= zobrist.SideToMove
	}

	// generate position
	ranks := strings.Split(parts[0], "/")
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
	board.CastlingRights = castling.NewRights(parts[2])
	board.Hash ^= zobrist.Castling[board.CastlingRights]

	// en-passant target square
	board.EnPassantTarget = square.NewFromString(parts[3])
	if board.EnPassantTarget != square.None {
		board.Hash ^= zobrist.EnPassant[board.EnPassantTarget.File()]
	}

	// move counters
	board.DrawClock, _ = strconv.Atoi(parts[4])
	board.FullMoves, _ = strconv.Atoi(parts[5])

	return &board
}

// FEN returns the fen string of the current Board position.
func (b *Board) FEN() string {
	var fenString string
	fenString += b.Position.FEN() + " "
	fenString += b.SideToMove.String() + " "
	fenString += b.CastlingRights.String() + " "
	fenString += b.EnPassantTarget.String() + " "
	fenString += strconv.Itoa(b.DrawClock) + " "
	fenString += strconv.Itoa(b.FullMoves)
	return fenString
}
