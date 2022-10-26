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

package move

import (
	"laptudirm.com/x/mess/pkg/chess/piece"
	"laptudirm.com/x/mess/pkg/chess/square"
)

type Move uint32

func New(source, target square.Square, from piece.Piece, isCapture bool) Move {
	m := Move(source) | Move(target<<6) | Move(from<<12) | Move(from<<16)
	if isCapture {
		m |= 0x100000
	}
	return m
}

func (m Move) String() string {
	s := m.Source().String() + m.Target().String()
	if m.IsPromotion() {
		s += m.ToPiece().String()
	}
	return s
}

func (m Move) SetPromotion(p piece.Piece) Move {
	m &^= 0xf0000
	m |= Move(p << 16)
	return m
}

func (m Move) Source() square.Square {
	return square.Square(m & 0x3f)
}

func (m Move) Target() square.Square {
	return square.Square((m & 0xfc0) >> 6)
}

func (m Move) FromPiece() piece.Piece {
	return piece.Piece((m & 0xf000) >> 12)
}

func (m Move) ToPiece() piece.Piece {
	return piece.Piece((m & 0xf0000) >> 16)
}

func (m Move) IsCapture() bool {
	return m&0x100000 != 0
}

func (m Move) IsPromotion() bool {
	return m.FromPiece() != m.ToPiece()
}

func (m Move) IsQuiet() bool {
	return !m.IsCapture() && !m.IsPromotion()
}

func (m Move) IsReversible() bool {
	return !m.IsCapture() && m.FromPiece().Type() != piece.Pawn
}

func (m Move) IsEnPassant(ep square.Square) bool {
	return m.Target() == ep && m.FromPiece().Type() == piece.Pawn
}
