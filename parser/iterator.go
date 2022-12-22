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
	// Will return false as the second argument if it cannot peek N times because
	// there are not as much items left
	peekN(n int) (*T, bool)
	// See the previous eleemnt
	prev() *T
	// Reverse the iterator back x amount of steps
	reverse(x int)
	// Take a subslices from this iterator starting at the current position (inclusive)
	// and until current position + n
	subslice(n int) Iterator[T]
	// Consumes n amount of items
	// Returns false if the new pointer exceeds the length of the array (e.g. more items where consumed than there are left)
	consume(n int) bool
	// Returns the amount of elements left in this iterator
	len() int
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

func (iter *ArrayIterator[T]) peekN(n int) (*T, bool) {
	if len(iter.arr) <= iter.ptr+n-1 || n == 0 { // Ptr is already at plus 1
		return nil, false
	}
	return &iter.arr[iter.ptr+n-1], true
}

func (iter *ArrayIterator[T]) prev() *T {
	return &iter.arr[iter.ptr-1]
}

func (iter *ArrayIterator[T]) reverse(x int) {
	iter.ptr -= x
}

func (iter *ArrayIterator[T]) subslice(n int) Iterator[T] {
	subslice := iter.arr[iter.ptr : iter.ptr+n]
	iterNew := newArrayIterator(subslice)
	return &iterNew
}

func (iter *ArrayIterator[T]) len() int {
	return len(iter.arr)
}

func (iter *ArrayIterator[T]) consume(n int) bool {
	iter.ptr += n
	return !(iter.ptr > len(iter.arr))
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

func (iter *ArrayOfPointersIterator[T]) peekN(n int) (*T, bool) {
	if len(iter.arr) <= iter.ptr+n-1 || n == 0 { // Ptr is already at plus 1
		return nil, false
	}
	return iter.arr[iter.ptr+n-1], true
}

func (iter *ArrayOfPointersIterator[T]) prev() *T {
	return iter.arr[iter.ptr-1]
}

func (iter *ArrayOfPointersIterator[T]) reverse(x int) {
	iter.ptr -= x
}

func (iter *ArrayOfPointersIterator[T]) subslice(n int) Iterator[T] {
	subslice := iter.arr[iter.ptr : iter.ptr+n]
	newIter := newArrayOfPointerIterator(subslice)
	return &newIter
}

func (iter *ArrayOfPointersIterator[T]) consume(n int) bool {
	iter.ptr += n
	return !(iter.ptr > len(iter.arr))
}

func (iter *ArrayOfPointersIterator[T]) len() int {
	return len(iter.arr)
}
