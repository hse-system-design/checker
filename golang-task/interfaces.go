package ordcol

import (
	"errors"
	"golang.org/x/exp/constraints"
)

var ErrDuplicateKey = errors.New("duplicate_key")
var ErrEmptyIterator = errors.New("empty_iterator")
var ErrEmptyCollection = errors.New("empty_collection")
var ErrUnknownOrder = errors.New("unknown_order")

type IterationOrder int

const (
	ByInsertion    IterationOrder = 1
	ByInsertionRev IterationOrder = 2
)

type Iterator[K comparable, V any] interface {
	HasNext() bool
	// Next returns ErrEmptyIterator if there is no more elements (HasNext() is false)
	Next() (K, V, error)
}

type Collection[K constraints.Ordered, V any] interface {
	// Add adds specified value to the collection and associates it with the specified key.
	//
	// If there already is a value associated with the same key, ErrDuplicateKey is returned.
	Add(key K, value V) error
	// DelMin removes the value associated with the minimum key from the collection and returns it.
	//
	// If the collection is empty, ErrEmptyCollection is returned.
	DelMin() (K, V, error)
	// Len returns the number of elements in the collection.
	Len() int

	// IterateBy returns an iterator over the elements of the collection in the specified order.
	// Any modification of the collection corrupts all the iterators obtained prior the modification.
	//
	// If invalid order is passed, the function panics with ErrUnknownOrder
	IterateBy(order IterationOrder) Iterator[K, V]
	// At returns the value associated with the specified key.
	At(key K) (V, bool)
}
