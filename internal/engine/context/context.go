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

package context

import (
	"laptudirm.com/x/mess/pkg/search"
	"laptudirm.com/x/mess/pkg/uci"
	"laptudirm.com/x/mess/pkg/uci/option"
)

// Engine represents the context containing the engine's information which
// is shared among it's UCI commands to store state.
type Engine struct {
	// engine's uci client
	Client uci.Client

	// current search context
	Search    *search.Context
	Searching bool

	Pondering    bool
	PonderLimits search.Limits

	// uci options
	OptionSchema option.Schema
	Options      options
}

// options contains the values of the UCI options supported by the engine.
type options struct {
	Ponder  bool // name Ponder type check
	Hash    int  // name Hash type spin
	Threads int  // name Threads type spin
}
