// Package analyzers for static code analyzers.
package analyzers

import "golang.org/x/tools/go/analysis"

// OsExitCheckAnalyzer a structure for an analyzer that checks for
// calling os.Exit in the main package, the main function
var OsExitCheckAnalyzer = &analysis.Analyzer{
	Name: "OsExitCheckAnalyzer",
	Doc:  "check for os.exit in main package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() == "main" {
		for _, file := range pass.Files {
			if file.Name.Name == "main" {
				for i, v := range pass.TypesInfo.Uses {
					if v.String() == "func os.Exit(code int)" {
						pass.Reportf(i.Pos(), "os.Exit in main file")
					}
				}
			}
		}
	}
	return nil, nil
}
