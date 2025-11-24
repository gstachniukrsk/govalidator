package govalidator

import (
	"context"
	"encoding/json"
)

// JSONPresenter creates a presenter that formats errors as JSON strings.
// Each error is converted to a JSON object with "path" and "message" fields.
//
// Example output: {"path":"$.user.age","message":"not an integer"}
func JSONPresenter(pathGlue string) PresenterFunc {
	pathPresenter := PathPresenter(pathGlue)

	return func(ctx context.Context, path []string, err error) string {
		pathStr := pathPresenter(ctx, path, err)

		errorData := map[string]string{
			"path":    pathStr,
			"message": err.Error(),
		}

		jsonBytes, jsonErr := json.Marshal(errorData)
		if jsonErr != nil {
			// fallback to simple format if JSON marshaling fails
			return pathStr + ": " + err.Error()
		}

		return string(jsonBytes)
	}
}

// JSONDetailedPresenter creates a presenter that formats errors as detailed JSON strings.
// It attempts to extract structured information from typed errors when available.
//
// Example output: {"path":"$.age","message":"not an integer","type":"NotAnIntegerError"}
func JSONDetailedPresenter(pathGlue string) PresenterFunc {
	pathPresenter := PathPresenter(pathGlue)

	return func(ctx context.Context, path []string, err error) string {
		pathStr := pathPresenter(ctx, path, err)

		errorData := map[string]any{
			"path":    pathStr,
			"message": err.Error(),
			"type":    getErrorType(err),
		}

		// Extract additional fields from structured errors
		switch e := err.(type) {
		case MinSizeError:
			errorData["minSize"] = e.MinSize
			errorData["actualSize"] = e.ActualSize
		case MaxSizeError:
			errorData["maxSize"] = e.MaxSize
			errorData["actualSize"] = e.ActualSize
		case FloatPrecisionError:
			errorData["expectedPrecision"] = e.ExpectedPrecision
			errorData["actualPrecision"] = e.ActualPrecision
		case FloatTooSmallError:
			errorData["minFloat"] = e.MinFloat
		case FloatTooLargeError:
			errorData["maxFloat"] = e.MaxFloat
		case StringTooShortError:
			errorData["minLength"] = e.MinLength
		case StringTooLongError:
			errorData["maxLength"] = e.MaxLength
		case FieldNotDefinedError:
			errorData["field"] = e.Field
		case UnexpectedFieldError:
			errorData["field"] = e.Field
		}

		jsonBytes, jsonErr := json.Marshal(errorData)
		if jsonErr != nil {
			// fallback to simple format if JSON marshaling fails
			return pathStr + ": " + err.Error()
		}

		return string(jsonBytes)
	}
}

// getErrorType returns the type name of the error
func getErrorType(err error) string {
	switch err.(type) {
	case RequiredError:
		return "RequiredError"
	case NotAStringError:
		return "NotAStringError"
	case NotAnIntegerError:
		return "NotAnIntegerError"
	case NotAFloatError:
		return "NotAFloatError"
	case NotABooleanError:
		return "NotABooleanError"
	case NotAMapError:
		return "NotAMapError"
	case NotAListError:
		return "NotAListError"
	case NotANumberError:
		return "NotANumberError"
	case MinSizeError:
		return "MinSizeError"
	case MaxSizeError:
		return "MaxSizeError"
	case FloatPrecisionError:
		return "FloatPrecisionError"
	case FloatTooSmallError:
		return "FloatTooSmallError"
	case FloatTooLargeError:
		return "FloatTooLargeError"
	case StringTooShortError:
		return "StringTooShortError"
	case StringTooLongError:
		return "StringTooLongError"
	case FieldNotDefinedError:
		return "FieldNotDefinedError"
	case UnexpectedFieldError:
		return "UnexpectedFieldError"
	default:
		return "Error"
	}
}
