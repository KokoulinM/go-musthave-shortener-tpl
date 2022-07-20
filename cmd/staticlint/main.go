// Package for statistical code analysis.
// To use, you need to compile/install.
// Run static int <path_to_files>, where path_to_files is the path to the files,
// which need to be checked.
// This package includes static pocket analyzers check.io,
// standard static package analyzers golang.org/x/tools/go/analysis/passes,
// OsExitCheckAnalyzer analyzer.
package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"staticlint/analyzers"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

const Config = `config.json`

type ConfigData struct {
	StaticCheck []string
}

func main() {
	appFile, err := os.Executable()
	if err != nil {
		panic(err)
	}

	data, err := os.ReadFile(filepath.Join(filepath.Dir(appFile), Config))
	if err != nil {
		panic(err)
	}

	var cfg ConfigData
	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	myChecks := []*analysis.Analyzer{
		analyzers.OsExitCheckAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		structtag.Analyzer,
	}

	checks := make(map[string]bool)

	for _, v := range cfg.StaticCheck {
		checks[v] = true
	}

	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			myChecks = append(myChecks, v.Analyzer)
		}
	}

	multichecker.Main(
		myChecks...,
	)
}
