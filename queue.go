package main

type QueueData interface {
	string | Signal
}

type QueueError struct {
	message string
}

func (e *QueueError) Error() string {
	return e.message
}

type Queue[T QueueData] struct {
	values []T
}

func (q *Queue[T]) Enqueue(value T) {
	q.values = append(q.values, value)
}

func (q *Queue[T]) Dequeue() T {
	if q.IsEmpty() {
		panic(&QueueError{message: "Cannot dequeue from an empty queue."})
	}

	dequeued := q.values[0]
	q.values = q.values[1:]

	return dequeued
}

func (q *Queue[T]) Peek() T {
	if q.IsEmpty() {
		panic(&QueueError{message: "Cannot peek an empty queue."})
	}

	return q.values[0]
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.values) == 0
}

func (q *Queue[T]) Size() int {
	return len(q.values)
}
