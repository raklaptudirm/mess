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

package search

import (
	"time"

	"laptudirm.com/x/mess/internal/util"
	"laptudirm.com/x/mess/pkg/board/piece"
)

// Manager represents a time manager.
type TimeManager interface {
	// GetDeadline calculates the optimal amount of time to be used
	// and sets a deadline internally for the search's end.
	GetDeadline()

	// ExtendDeadline is called when the engine want's to extend the
	// search's length. A deadline extension may fail.
	ExtendDeadline()

	// PessimisticExpired reports if the search deadline has been crossed.
	PessimisticExpired() bool
	OptimisticExpired() bool
}

// NormalManager is the standard time manager which uses the wtime, btime,
// winc, binc, and movestogo provided by the GUI to calculate the optimal
// search time.
type TimeManagerNormal struct {
	Us piece.Color // side to move

	Time, Increment [piece.ColorN]int
	MovesToGo       int // moves to next time control

	deadline time.Time // search end deadline
	maxUsage time.Time

	context *Context
}

// compile time check that NormalManager implements Manager
var _ TimeManager = (*TimeManagerNormal)(nil)

func (c *TimeManagerNormal) GetDeadline() {
	if c.MovesToGo == 0 {
		c.MovesToGo = util.Max(20, 50-(c.context.board.Plys/2))
	}

	maximum := time.Duration(c.Time[c.Us]) * time.Millisecond

	c.deadline = time.Now().Add(maximum / time.Duration(c.MovesToGo))
	c.maxUsage = time.Now().Add(maximum / 2)
}

func (c *TimeManagerNormal) ExtendDeadline() {
	c.deadline = c.deadline.Add((time.Duration(c.Time[c.Us]) * time.Millisecond) / 30)
}

func (c *TimeManagerNormal) PessimisticExpired() bool {
	return time.Now().After(c.maxUsage)
}

func (c *TimeManagerNormal) OptimisticExpired() bool {
	return time.Now().After(c.deadline)
}

// MoveManger is the time manager used when the gui wants to time a search
// by move-time. Extending it's deadline is not possible.
type TimeManagerMovetime struct {
	deadline time.Time
	Duration int
}

// compile time check that MoveManager implements Manager
var _ TimeManager = (*TimeManagerMovetime)(nil)

func (c *TimeManagerMovetime) GetDeadline() {
	c.deadline = time.Now().Add(time.Duration(c.Duration) * time.Millisecond)
}

func (c *TimeManagerMovetime) ExtendDeadline() {
	// can't extend deadline: search time is fixed
}

func (c *TimeManagerMovetime) PessimisticExpired() bool {
	return time.Now().After(c.deadline)
}

func (c *TimeManagerMovetime) OptimisticExpired() bool {
	return time.Now().After(c.deadline)
}
