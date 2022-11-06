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

// Package generator implements utility function used by code generators.
package generator

import (
	"os"
	"text/template"
)

// Generate evaluates the given template string t with the data v and
// writes it to a new generated file with the name <name>.go in the cwd.
func Generate(name, t string, v any) {
	f, err := os.Create(name + ".go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	template.Must(template.New(name).Parse(t)).Execute(f, v)
}
