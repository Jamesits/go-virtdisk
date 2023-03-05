package struct_alignment

// from: golang.org/x/tools@v0.5.0/go/analysis/singlechecker/singlechecker.go

import (
	"fmt"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
	"sort"
	"strings"
)

func Run(patterns []string) int {
	analyzers := []*analysis.Analyzer{
		fieldalignment.Analyzer,
	}

	if err := analysis.Validate(analyzers); err != nil {
		log.Fatal(err)
	}

	// Optimization: if the selected analyzers don't produce/consume
	// facts, we need source only for the initial packages.
	allSyntax := needFacts(analyzers)
	initial, err := load(patterns, allSyntax)
	if err != nil {
		if _, ok := err.(typeParseError); !ok {
			// Fail when some of the errors are not
			// related to parsing nor typing.
			log.Print(err)
			return 1
		}
	}

	// Print the results.
	roots := analyze(initial, analyzers)
	return printDiagnostics(roots)
}

// exportedFrom reports whether obj may be visible to a package that imports pkg.
// This includes not just the exported members of pkg, but also unexported
// constants, types, fields, and methods, perhaps belonging to other packages,
// that find there way into the API.
// This is an overapproximation of the more accurate approach used by
// gc export data, which walks the type graph, but it's much simpler.
//
// TODO(adonovan): do more accurate filtering by walking the type graph.
func exportedFrom(obj types.Object, pkg *types.Package) bool {
	switch obj := obj.(type) {
	case *types.Func:
		return obj.Exported() && obj.Pkg() == pkg ||
			obj.Type().(*types.Signature).Recv() != nil
	case *types.Var:
		if obj.IsField() {
			return true
		}
		// we can't filter more aggressively than this because we need
		// to consider function parameters exported, but have no way
		// of telling apart function parameters from local variables.
		return obj.Pkg() == pkg
	case *types.TypeName, *types.Const:
		return true
	}
	return false // Nil, Builtin, Label, or PkgName
}

// load loads the initial packages. If all loading issues are related to
// typing and parsing, the returned error is of type typeParseError.
func load(patterns []string, allSyntax bool) ([]*packages.Package, error) {
	mode := packages.LoadSyntax
	if allSyntax {
		mode = packages.LoadAllSyntax
	}
	conf := packages.Config{
		Mode:  mode,
		Tests: true,
	}
	initial, err := packages.Load(&conf, patterns...)
	if err == nil {
		if len(initial) == 0 {
			err = fmt.Errorf("%s matched no packages", strings.Join(patterns, " "))
		} else {
			err = loadingError(initial)
		}
	}
	return initial, err
}

func analyze(pkgs []*packages.Package, analyzers []*analysis.Analyzer) []*action {
	// Construct the action graph.

	// Each graph node (action) is one unit of analysis.
	// Edges express package-to-package (vertical) dependencies,
	// and analysis-to-analysis (horizontal) dependencies.
	type key struct {
		*analysis.Analyzer
		*packages.Package
	}
	actions := make(map[key]*action)

	var mkAction func(a *analysis.Analyzer, pkg *packages.Package) *action
	mkAction = func(a *analysis.Analyzer, pkg *packages.Package) *action {
		k := key{a, pkg}
		act, ok := actions[k]
		if !ok {
			act = &action{a: a, pkg: pkg}

			// Add a dependency on each required analyzers.
			for _, req := range a.Requires {
				act.deps = append(act.deps, mkAction(req, pkg))
			}

			// An analysis that consumes/produces facts
			// must run on the package's dependencies too.
			if len(a.FactTypes) > 0 {
				paths := make([]string, 0, len(pkg.Imports))
				for path := range pkg.Imports {
					paths = append(paths, path)
				}
				sort.Strings(paths) // for determinism
				for _, path := range paths {
					dep := mkAction(a, pkg.Imports[path])
					act.deps = append(act.deps, dep)
				}
			}

			actions[k] = act
		}
		return act
	}

	// Build nodes for initial packages.
	var roots []*action
	for _, a := range analyzers {
		for _, pkg := range pkgs {
			root := mkAction(a, pkg)
			root.isroot = true
			roots = append(roots, root)
		}
	}

	// Execute the graph in parallel.
	execAll(roots)

	return roots
}

func execAll(actions []*action) {
	for _, act := range actions {
		act.exec()
	}
}

// printDiagnostics prints the diagnostics for the root packages in either
// plain text or JSON format. JSON format also includes errors for any
// dependencies.
//
// It returns the exitcode: in plain mode, 0 for success, 1 for analysis
// errors, and 3 for diagnostics. We avoid 2 since the flag package uses
// it. JSON mode always succeeds at printing errors and diagnostics in a
// structured form to stdout.
func printDiagnostics(roots []*action) (exitcode int) {
	// Print the output.
	//
	// Print diagnostics only for root packages,
	// but errors for all packages.
	printed := make(map[*action]bool)
	var print func(*action)
	var visitAll func(actions []*action)
	visitAll = func(actions []*action) {
		for _, act := range actions {
			if !printed[act] {
				printed[act] = true
				visitAll(act.deps)
				print(act)
			}
		}
	}

	// plain text output

	// De-duplicate diagnostics by position (not token.Pos) to
	// avoid double-reporting in source files that belong to
	// multiple packages, such as foo and foo.test.
	type key struct {
		pos token.Position
		end token.Position
		*analysis.Analyzer
		message string
	}
	seen := make(map[key]bool)

	print = func(act *action) {
		if act.err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", act.a.Name, act.err)
			exitcode = 1 // analysis failed, at least partially
			return
		}
		if act.isroot {
			for _, diag := range act.diagnostics {
				// We don't display a.Name/f.Category
				// as most users don't care.

				posn := act.pkg.Fset.Position(diag.Pos)
				end := act.pkg.Fset.Position(diag.End)
				k := key{posn, end, act.a, diag.Message}
				if seen[k] {
					continue // duplicate
				}
				seen[k] = true

				PrintPlain(act.pkg.Fset, diag, 0)
			}
		}
	}
	visitAll(roots)

	if exitcode == 0 && len(seen) > 0 {
		exitcode = 3 // successfully produced diagnostics
	}

	return exitcode
}
