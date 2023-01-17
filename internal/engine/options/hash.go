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

package options

import (
	"laptudirm.com/x/mess/internal/engine/context"
	"laptudirm.com/x/mess/pkg/uci/option"
)

// UCI option Hash, type spin
//
// The value in MB allocated for hash tables.
// This should be answered with the first "setoption" command at program
// boot if the engine has sent the appropriate option name Hash command,
// which should be supported by all engines! So the engine should use a
// very small hash value as default.
func NewHash(engine *context.Engine) option.Option {
	return &option.Spin{
		Default: 16, // default from stockfish
		Min:     1,
		// use stockfish value to suppress cutechess warnings
		Max: 33554432,
		Storage: func(hash int) error {
			engine.Options.Hash = hash

			// resize hash table
			engine.Search.ResizeTT(hash)

			return nil
		},
	}
}
