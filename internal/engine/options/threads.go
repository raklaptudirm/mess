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

// UCI option Threads, type spin
//
// The number of threads the engine should use while searching.
func NewThreads(engine *context.Engine) option.Option {
	return &option.Spin{
		// multi-threading not implemented, so fix value at 1
		Default: 1,
		Min:     1, Max: 1,

		Storage: func(threads int) error {
			engine.Options.Threads = threads
			return nil
		},
	}
}
