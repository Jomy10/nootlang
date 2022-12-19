package parser

// A generic iterator over elements of type `T`
type Iterator[T any] interface {
	// Get the next element of the iterator. Returns false as the second argument
	// if there are no more elements left in the Iterator
	next() (*T, bool)
	// Returns true if there is another element in the iterator
	hasNext() bool
	// Returns the next element of the iterator without advancing it. Returns false
	// as the second argument if there are no more elements in the iterator
	peek() (*T, bool)
}

// A simple iterator over an array of values T
type ArrayIterator[T any] struct {
	arr []T
	// points to the next item in the iterator
	ptr int
}

func newArrayIterator[T any](arr []T) ArrayIterator[T] {
	return ArrayIterator[T]{arr, 0}
}

// Returns the next elements of the iterator and true, or false if there are no
// more elements in the iterator
func (iter *ArrayIterator[T]) next() (*T, bool) {
	if !iter.hasNext() {
		return nil, false
	}

	val := iter.arr[iter.ptr]
	iter.ptr++
	return &val, true
}

// Returns wheter the iterator has more elements
func (iter *ArrayIterator[T]) hasNext() bool {
	return len(iter.arr) != iter.ptr
}

// Peek at the next eleemnt of the iterator, without consuming it
func (iter *ArrayIterator[T]) peek() (*T, bool) {
	if !iter.hasNext() {
		return nil, false
	}
	return &iter.arr[iter.ptr], true
}

type ArrayOfPointersIterator[T any] struct {
	arr []*T
	ptr int
}

func newArrayOfPointerIterator[T any](arr []*T) ArrayOfPointersIterator[T] {
	return ArrayOfPointersIterator[T]{arr, 0}
}

func (iter *ArrayOfPointersIterator[T]) next() (*T, bool) {
	if !iter.hasNext() {
		return nil, false
	}

	val := iter.arr[iter.ptr]
	iter.ptr++
	return val, true
}

func (iter *ArrayOfPointersIterator[T]) hasNext() bool {
	return len(iter.arr) != iter.ptr
}

func (iter *ArrayOfPointersIterator[T]) peek() (*T, bool) {
	if !iter.hasNext() {
		return nil, false
	}
	return iter.arr[iter.ptr], true
}
