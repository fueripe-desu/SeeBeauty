package main

type StackData interface {
	string | Signal
}

type StackError struct {
	message string
}

func (e *StackError) Error() string {
	return e.message
}

type Stack[T StackData] struct {
	values []T
}

func (s *Stack[T]) Push(value T) {
	s.values = append(s.values, value)
}

func (s *Stack[T]) Pop() T {
	if s.IsEmpty() {
		panic(&StackError{message: "Cannot pop if stack is empty."})
	}

	popped := s.Peek()
	s.values = s.values[:s.Size()-1]

	return popped
}

func (s *Stack[T]) Peek() T {
	if s.IsEmpty() {
		panic(&StackError{message: "Cannot peek empty stack."})
	}

	return s.values[s.Size()-1]
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.values) == 0
}

func (s *Stack[T]) Size() int {
	return len(s.values)
}
