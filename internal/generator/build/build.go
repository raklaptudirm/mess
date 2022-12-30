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

// This is a generator package used to generate go files containing Mess's
// build information, like version, etc.
package main

import (
	_ "embed"
	"os"
	"os/exec"
	"strings"

	"laptudirm.com/x/mess/internal/generator"
)

type buildStruct struct {
	Version string
}

//go:embed .gotemplate
var template string

func main() {
	var b buildStruct

	var err error
	b.Version, err = runOutput("git", "describe", "--tags")
	if err != nil {
		b.Version, err = runOutput("git", "rev-parse", "--short", "HEAD")
		if err != nil {
			b.Version = "v0.0.0"
		}
	}

	generator.Generate("info", template, b)
}

func runOutput(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stderr = os.Stderr   // print any errors
	out, err := cmd.Output() // copy the stdout

	return strings.TrimSuffix(string(out), "\n"), err
}
