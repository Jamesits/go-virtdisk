package struct_alignment

// from: golang.org/x/tools@v0.5.0/go/analysis/internal/analysisflags/flags.go

import (
	"fmt"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"io/ioutil"
	"os"
	"strings"
)

// PrintPlain prints a diagnostic in plain text form,
// with context specified by the -c flag.
// -c=N: if N>0, display offending line plus N lines of context
func PrintPlain(fset *token.FileSet, diag analysis.Diagnostic, Context int) {
	posn := fset.Position(diag.Pos)
	fmt.Fprintf(os.Stderr, "%s: %s\n", posn, diag.Message)

	// -c=N: show offending line plus N lines of context.
	if Context >= 0 {
		posn := fset.Position(diag.Pos)
		end := fset.Position(diag.End)
		if !end.IsValid() {
			end = posn
		}
		data, _ := ioutil.ReadFile(posn.Filename)
		lines := strings.Split(string(data), "\n")
		for i := posn.Line - Context; i <= end.Line+Context; i++ {
			if 1 <= i && i <= len(lines) {
				fmt.Fprintf(os.Stderr, "%d\t%s\n", i, lines[i-1])
			}
		}
	}
}
