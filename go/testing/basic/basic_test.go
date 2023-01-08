package main

import "testing"

// Run basic tests by
// $ go test
// Or printing results in verbose mode
// $ go test -v
// Generate coverage result by
// $ go test --cover
// And generate reports in HTML
// $ go test --cover --coverprofile .data/coverage.out
// $ go tool cover --html .data/coverage.out -o .data/coverage.html

func TestSumCorrect(t *testing.T) {
	testCases := []struct {
		a    int
		b    int
		want int
	}{
		{1, 2, 3},
		{3, 5, 8},
		{-1, 9, 8},
	}

	for _, testCase := range testCases {
		got := SumCorrect(testCase.a, testCase.b)
		if got != testCase.want {
			t.Errorf("Sum of %d and %d is incorrect: got %d, want %d",
				testCase.a, testCase.b, got, testCase.want)
		}
	}
}

func TestSumBuggy(t *testing.T) {
	testCases := []struct {
		a    int
		b    int
		want int
	}{
		{1, 2, 3},
		{3, 5, 8},
		{-1, 9, 8},
	}

	for _, testCase := range testCases {
		got := SumBuggy(testCase.a, testCase.b)
		if got != testCase.want {
			t.Errorf("Sum of %d and %d is incorrect: got %d, want %d",
				testCase.a, testCase.b, got, testCase.want)
		}
	}
}
