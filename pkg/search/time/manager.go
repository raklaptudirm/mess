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

// Package time implements various types and functions used to manage
// search time while searching a position.
package time

import (
	"time"

	"laptudirm.com/x/mess/pkg/board/piece"
)

// Manager represents a time manager.
type Manager interface {
	// GetDeadline calculates the optimal amount of time to be used
	// and sets a deadline internally for the search's end.
	GetDeadline()

	// ExtendDeadline is called when the engine want's to extend the
	// search's length. A deadline extension may fail.
	ExtendDeadline()

	// Expired reports if the search deadline has been crossed.
	Expired() bool
}

// NormalManager is the standard time manager which uses the wtime, btime,
// winc, binc, and movestogo provided by the GUI to calculate the optimal
// search time.
type NormalManager struct {
	Us piece.Color // side to move

	Time, Increment [piece.ColorN]int
	MovesToGo       int // moves to next time control

	deadline time.Time // search end deadline
}

// compile time check that NormalManager implements Manager
var _ Manager = (*NormalManager)(nil)

func (c *NormalManager) GetDeadline() {
	c.deadline = time.Now().Add((time.Duration(c.Time[c.Us]) * time.Millisecond) / 20)
}

func (c *NormalManager) ExtendDeadline() {
	c.deadline = c.deadline.Add((time.Duration(c.Time[c.Us]) * time.Millisecond) / 30)
}

func (c *NormalManager) Expired() bool {
	return time.Now().After(c.deadline)
}

// MoveManger is the time manager used when the gui wants to time a search
// by move-time. Extending it's deadline is not possible.
type MoveManager struct {
	Duration int
	deadline time.Time
}

// compile time check that MoveManager implements Manager
var _ Manager = (*MoveManager)(nil)

func (c *MoveManager) GetDeadline() {
	c.deadline = time.Now().Add(time.Duration(c.Duration) * time.Millisecond)
}

func (c *MoveManager) ExtendDeadline() {
	// can't extend deadline: search time is fixed
}

func (c *MoveManager) Expired() bool {
	return time.Now().After(c.deadline)
}
