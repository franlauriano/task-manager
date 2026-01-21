package assert

import (
	"reflect"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// CompareErrors compares two errors generically, detecting real content differences.
func CompareErrors(got, want error) string {
	if got == nil && want == nil {
		return ""
	}
	if got == nil || want == nil {
		return cmp.Diff(got, want, cmpopts.EquateErrors())
	}

	gotType := reflect.TypeOf(got)
	wantType := reflect.TypeOf(want)
	if gotType != wantType {
		// If types differ, compare error messages first, then types
		if got.Error() != want.Error() {
			return cmp.Diff(got.Error(), want.Error())
		}
		return cmp.Diff(gotType.String(), wantType.String())
	}

	diffWithEquate := cmp.Diff(got, want, cmpopts.EquateErrors())
	if diffWithEquate == "" {
		return ""
	}

	// For pointer-to-struct errors, try deep comparison without EquateErrors
	// Use recover to handle potential panics from unexported fields
	if gotType.Kind() == reflect.Ptr && gotType.Elem().Kind() == reflect.Struct {
		var diff string
		func() {
			defer func() {
				recover()
			}()
			diff = cmp.Diff(got, want)
		}()
		if diff != "" {
			return diff
		}
		return ""
	}

	// If error messages match, consider them equal
	if got.Error() == want.Error() {
		return ""
	}

	return diffWithEquate
}
