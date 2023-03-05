package struct_alignment

import (
	"errors"
	"fmt"
	"golang.org/x/tools/go/packages"
)

// typeParseError represents a package load error
// that is related to typing and parsing.
type typeParseError struct {
	error
}

// loadingError checks for issues during the loading of initial
// packages. Returns nil if there are no issues. Returns error
// of type typeParseError if all errors, including those in
// dependencies, are related to typing or parsing. Otherwise,
// a plain error is returned with an appropriate message.
func loadingError(initial []*packages.Package) error {
	var err error
	if n := packages.PrintErrors(initial); n > 1 {
		err = fmt.Errorf("%d errors during loading", n)
	} else if n == 1 {
		err = errors.New("error during loading")
	} else {
		// no errors
		return nil
	}
	all := true
	packages.Visit(initial, nil, func(pkg *packages.Package) {
		for _, err := range pkg.Errors {
			typeOrParse := err.Kind == packages.TypeError || err.Kind == packages.ParseError
			all = all && typeOrParse
		}
	})
	if all {
		return typeParseError{err}
	}
	return err
}
