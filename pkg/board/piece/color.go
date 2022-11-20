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

package piece

// NewColor creates an instance of color from the given id.
func NewColor(id string) Color {
	switch id {
	case "w":
		return White
	case "b":
		return Black
	default:
		panic("new color: invalid color id")
	}
}

// Color represents the color of a Piece.
type Color int

// constants representing various piece colors
const (
	White Color = iota
	Black
)

// ColorN is the number of colors there are.
const ColorN = 2

// Other returns the color opposite to the given one. For White,
// it returns Black and for Black, it returns White.
func (c Color) Other() Color {
	return 1 ^ c
}

// String converts a Color to it's string representation.
func (c Color) String() string {
	const colorToStr = "wb"
	return string(colorToStr[c])
}
