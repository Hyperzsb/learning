package parallel

import "time"

var addon = 0

func AddCorrect(a int) int {
	return a + 1
}

func AddBuggy(a int) int {
	addon = a
	time.Sleep(time.Second)
	addon = 1
	defer func() {
		addon = 0
	}()

	return a + addon
}
