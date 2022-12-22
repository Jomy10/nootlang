package parser

// A simple iterator over an array of values T
type Iterator[T any] struct {
	arr []T
	// Points to the next item in the iterator
	ptr int
}

func newIterator[T any](arr []T) Iterator[T] {
	return Iterator[T]{
		arr,
		0,
	}
}

// Returns the next elements of the iterator and true, or false if there are no
// more elements in the iterator
func (iter *Iterator[T]) next() (*T, bool) {
	if !iter.hasNext() {
		return nil, false
	}

	val := iter.arr[iter.ptr]
	iter.ptr++
	return &val, true
}

// Returns wheter thee iterator has more elements
func (iter *Iterator[T]) hasNext() bool {
	return len(iter.arr) != iter.ptr
}

// Peek at the next eleemnt of the iterator, without consuming it
func (iter *Iterator[T]) peek() *T {
	if !iter.hasNext() {
		return nil
	}
	return &iter.arr[iter.ptr]
}
