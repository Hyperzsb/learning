package parallel

import "testing"

// Run basic tests in parallel by
// $ go test --parallel 5

func TestAddCorrect(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		a    int
		want int
	}{
		{"1", 1, 2},
		{"2", -1, 0},
		{"3", 99, 100},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := AddCorrect(tc.a)
			if got != tc.want {
				t.Errorf("#%s failed: want %d got %d", tc.name, tc.want, got)
			}
		})
	}
}

func TestAddBuggy(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		a    int
		want int
	}{
		{"1", 1, 2},
		{"2", -1, 0},
		{"3", 99, 100},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := AddBuggy(tc.a)
			if got != tc.want {
				t.Errorf("#%s failed: want %d got %d", tc.name, tc.want, got)
			}
		})
	}
}
