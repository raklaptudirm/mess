package main

import (
	"fmt"
	"os"
	"strings"

	"laptudirm.com/x/mess/scripts/util"
)

func main() {
	var args []string

	// set env variables from args
	for _, arg := range os.Args[1:] {
		name, value, found := strings.Cut(arg, "=")
		if !found {
			args = append(args, arg)
			continue
		}

		// var=value
		os.Setenv(name, value)
	}

	// remaining args are tasks
	for _, arg := range args {
		task, ok := tasks[arg]
		if !ok {
			fmt.Fprintf(os.Stderr, "Invalid task %v.\n", arg)
			continue
		}

		if err := task(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

var tasks = map[string]func() error{
	"--": nullTask, // used as separator for readability

	"dev-build":     devBuild,     // build a development binary
	"release-build": releaseBuild, // build a release binary
}

func nullTask() error {
	return nil
}

func devBuild() error {
	// version is latest tag-commits after tag-current commit hash
	version, err := util.RunWithOutput("git", "describe", "--tags")
	if err != nil {
		return err
	}

	return build(version)
}

func releaseBuild() error {
	// version is latest tag
	version, err := util.RunWithOutput("git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		return err
	}

	return build(version)
}

func build(version string) error {
	project := "laptudirm.com/x/mess"
	ldflags := fmt.Sprintf("-X %s/internal/build.Version=%s", project, version)

	var exe string
	if exe = os.Getenv("EXE"); exe == "" {
		exe = "mess"
	}

	if os.Getenv("GOOS") == "windows" {
		exe += ".exe"
	}

	return util.RunNormal("go", "build", "-ldflags", ldflags, "-o", exe)
}
