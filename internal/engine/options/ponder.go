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

// UCI option Ponder, type check
//
// This means that the engine is able to ponder. The GUI will send this
// whenever pondering is possible or not.
//
// Note: The engine should not start pondering on its own if this is
// enabled, this option is only needed because the engine might change its
// time management algorithm when pondering is allowed.
func NewPonder(engine *context.Engine) option.Option {
	return &option.Check{
		Default: false,
		Storage: func(ponder bool) error {
			engine.Options.Ponder = ponder
			return nil
		},
	}
}
