package govalidator

import (
	"context"
	"strings"
)

// CombinedPresenter creates a presenter that combines path and error message into a single string.
// It uses pathGlue to join path segments (e.g., ".") and separator to separate path from error (e.g., ": ").
//
// Example output: "$.user.age: not an integer"
func CombinedPresenter(pathGlue, separator string) PresenterFunc {
	pathPresenter := PathPresenter(pathGlue)
	errorPresenter := SimpleErrorPresenter()

	return func(ctx context.Context, path []string, err error) string {
		pathStr := pathPresenter(ctx, path, err)
		errorStr := errorPresenter(ctx, path, err)

		if pathStr == "" {
			return errorStr
		}

		return pathStr + separator + errorStr
	}
}

// CombinedBracketPresenter creates a presenter that combines path and error using bracket notation.
// Similar to CombinedPresenter but formats the path with brackets for arrays.
//
// Example output: "$[0].user.age: not an integer"
func CombinedBracketPresenter(pathGlue, separator string) PresenterFunc {
	return func(ctx context.Context, path []string, err error) string {
		if len(path) == 0 {
			return err.Error()
		}

		// build path with smart bracket handling
		var result strings.Builder
		for i, p := range path {
			if i == 0 {
				result.WriteString(p)
				continue
			}

			// check if this segment is an array index
			if strings.HasPrefix(p, "[") && strings.HasSuffix(p, "]") {
				result.WriteString(p)
			} else {
				result.WriteString(pathGlue)
				result.WriteString(p)
			}
		}

		return result.String() + separator + err.Error()
	}
}
