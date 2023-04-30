package queue

type Queue[T any] []T

func (q *Queue[T]) Push(n T) {
	*q = append(*q, n)
}

func (q *Queue[T]) Pop() {
	*q = (*q)[1:]
}

func (q *Queue[T]) Front() T {
	return (*q)[0]
}

func (q *Queue[T]) Empty() bool {
	return len(*q) == 0
}

func (q *Queue[T]) Length() int {
	return len(*q)
}
