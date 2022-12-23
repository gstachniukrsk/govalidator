package main_test

import "testing"

func assertContainsAllErrors(t *testing.T, haystack []error, needles []error) {
	for _, needle := range needles {
		assertContainsError(t, haystack, needle)
	}
}

func assertContainsError(t *testing.T, haystack []error, err error) {
	for _, e := range haystack {
		if e == err {
			return
		}
	}
	t.Errorf("expected %v to contain %v", haystack, err)
}

func anyPtr(any any) any {
	return &any
}

func strPtr(str string) *string {
	return &str
}
