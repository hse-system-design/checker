package ordcol

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
	"math/rand"
	"testing"
)

func TestCollection_Correctness(t *testing.T) {
	t.Run("AddAndGet", func(t *testing.T) {
		col := NewCollection[int, int]()

		err := col.Add(2, 4)
		require.NoError(t, err)
		err = col.Add(5, 9)
		require.NoError(t, err)

		t.Run("PresentKey", func(t *testing.T) {
			val, ok := col.At(2)
			require.True(t, ok)
			require.Equal(t, 4, val)

			val, ok = col.At(5)
			require.True(t, ok)
			require.Equal(t, 9, val)
		})

		t.Run("AbsentKey", func(t *testing.T) {
			_, ok := col.At(42)
			require.False(t, ok)
		})

		t.Run("AddExistingKey", func(t *testing.T) {
			err := col.Add(2, 5)
			require.ErrorIs(t, err, ErrDuplicateKey)
		})

	})

	t.Run("Empty", func(t *testing.T) {
		col := NewCollection[int, int]()

		t.Run("At", func(t *testing.T) {
			_, ok := col.At(1)
			require.False(t, ok)
		})

		t.Run("IterByInsertion", func(t *testing.T) {
			it := col.IterateBy(ByInsertion)
			require.False(t, it.HasNext())
			_, _, err := it.Next()
			require.ErrorIs(t, err, ErrEmptyIterator)
		})

		t.Run("IterByInsertionRev", func(t *testing.T) {
			it := col.IterateBy(ByInsertionRev)
			require.False(t, it.HasNext())
			_, _, err := it.Next()
			require.ErrorIs(t, err, ErrEmptyIterator)
		})

		t.Run("InvalidOrder", func(t *testing.T) {
			var rawErr any
			func() {
				defer func() { rawErr = recover() }()
				col.IterateBy(100500)
			}()
			err, ok := rawErr.(error)
			require.True(t, ok, "expected erorr but got %v", err)
			require.ErrorIs(t, err, ErrUnknownOrder)
		})
	})

	t.Run("IterateByInsertion", func(t *testing.T) {
		keys := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		rand.Shuffle(len(keys), func(i, j int) {
			keys[i], keys[j] = keys[j], keys[i]
		})

		col := NewCollection[int, int]()
		for _, key := range keys {
			err := col.Add(key, key*2)
			require.NoError(t, err)
		}
		require.Equal(t, len(keys), col.Len())

		t.Run("AscOrder", func(t *testing.T) {
			var expectedOrder []int
			for _, key := range keys {
				expectedOrder = append(expectedOrder, key, key*2)
			}

			order := drainIterator(t, col.IterateBy(ByInsertion))
			require.Equal(t, expectedOrder, order)
		})

		t.Run("DescOrder", func(t *testing.T) {
			expectedOrder := make([]int, 2*len(keys))
			for idx, key := range keys {
				expectedOrder[len(expectedOrder)-2*idx-2] = key
				expectedOrder[len(expectedOrder)-2*idx-1] = 2 * key
			}

			order := drainIterator(t, col.IterateBy(ByInsertionRev))
			require.Equal(t, expectedOrder, order)
		})
	})

	t.Run("InvalidateIterator", func(t *testing.T) {
		keys := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		rand.Shuffle(len(keys), func(i, j int) {
			keys[i], keys[j] = keys[j], keys[i]
		})

		col := NewCollection[int, int]()
		for _, key := range keys {
			err := col.Add(key, key*2)
			require.NoError(t, err)
		}
		require.Equal(t, len(keys), col.Len())

		limit := 5

		t.Run("AscOrder", func(t *testing.T) {
			var expectedOrder []int
			for _, key := range keys[:limit] {
				expectedOrder = append(expectedOrder, key)
			}

			addFn := func() error {
				return col.Add(0, 0)
			}
			delFn := func() error {
				key, value, err := col.DelMin()
				require.Equal(t, 0, key)
				require.Equal(t, 0, value)
				return err
			}

			for _, cb := range []func() error{addFn, delFn} {
				it := col.IterateBy(ByInsertion)
				order := drainKeysLimit(t, it, limit)
				require.Equal(t, expectedOrder, order)

				require.NoError(t, cb())
				_, _, err := it.Next()
				require.ErrorIs(t, err, ErrEmptyIterator)
			}
		})

		t.Run("DescOrder", func(t *testing.T) {
			expectedOrder := make([]int, limit)
			for idx, key := range keys[len(keys)-limit:] {
				expectedOrder[limit-idx-1] = key
			}

			addFn := func() error {
				copy(expectedOrder[1:limit], expectedOrder[:limit-1])
				expectedOrder[0] = 0
				return col.Add(0, 0)
			}
			delFn := func() error {
				key, value, err := col.DelMin()
				require.Equal(t, 0, key)
				require.Equal(t, 0, value)
				return err
			}

			for _, cb := range []func() error{addFn, delFn} {
				it := col.IterateBy(ByInsertionRev)
				order := drainKeysLimit(t, it, limit)
				require.Equal(t, expectedOrder, order)

				require.NoError(t, cb())
				_, _, err := it.Next()
				require.ErrorIs(t, err, ErrEmptyIterator)
			}
		})
	})

	t.Run("DelMin", func(t *testing.T) {
		orderedKeys := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		keys := slices.Clone(orderedKeys)
		rand.Shuffle(len(keys), func(i, j int) {
			keys[i], keys[j] = keys[j], keys[i]
		})

		col := NewCollection[int, int]()
		for _, key := range keys {
			err := col.Add(key, key*3)
			require.NoError(t, err)
		}
		require.Equal(t, len(keys), col.Len())

		var deletedKeys []int
		t.Run("DeleteFirstHalf", func(t *testing.T) {
			for i := 0; i < len(orderedKeys)/2; i++ {
				key, val, err := col.DelMin()
				require.NoError(t, err)
				require.Equal(t, key*3, val)
				deletedKeys = append(deletedKeys, key)
			}

			require.Equal(t, orderedKeys[:len(deletedKeys)], deletedKeys)
			require.ElementsMatch(t, orderedKeys[len(deletedKeys):], drainKeys(t, col.IterateBy(ByInsertion)))
		})

		t.Run("DeleteRemaining", func(t *testing.T) {
			for col.Len() > 0 {
				key, val, err := col.DelMin()
				require.NoError(t, err)
				require.Equal(t, key*3, val)
				deletedKeys = append(deletedKeys, key)
			}
			_, _, err := col.DelMin()
			require.ErrorIs(t, err, ErrEmptyCollection)

			require.Equal(t, orderedKeys, deletedKeys)
		})
	})
}

func TestCollection_Stress(t *testing.T) {
	numbers := make([]int, 1_234_567)
	for i := range numbers {
		numbers[i] = i + 1
	}
	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	col := NewCollection[int, int]()
	// insert all numbers
	for idx, num := range numbers {
		if idx > 0 && idx%10_000 == 0 {
			keys := drainKeysLimit(t, col.IterateBy(ByInsertionRev), 10_000)
			expected := slices.Clone(numbers[idx-10_000 : idx])
			for i := 0; i < len(expected)/2; i++ {
				j := len(expected) - i - 1
				expected[i], expected[j] = expected[j], expected[i]
			}
			require.Equal(t, expected, keys)
		}

		err := col.Add(num, num*5)
		require.NoError(t, err)

		require.Equal(t, num, drainKeysLimit(t, col.IterateBy(ByInsertionRev), 1)[0])
	}

	var deletedCount, earliestInsertedIdx int
	for col.Len() > 0 {
		key, val, err := col.DelMin()
		require.NoError(t, err)
		deletedCount++
		require.Equal(t, deletedCount, key)
		require.Equal(t, key*5, val)

		for earliestInsertedIdx < len(numbers) && numbers[earliestInsertedIdx] <= key {
			earliestInsertedIdx++
		}
		if col.Len() > 0 {
			earliestInserted := drainKeysLimit(t, col.IterateBy(ByInsertion), 1)[0]
			require.Equal(t, numbers[earliestInsertedIdx], earliestInserted)
		}
	}
}

func drainKeysLimit(t *testing.T, it Iterator[int, int], limit int) []int {
	var result []int
	for len(result) < limit && it.HasNext() {
		k, _, err := it.Next()
		require.NoError(t, err)
		result = append(result, k)
	}
	return result
}

func drainKeys(t *testing.T, it Iterator[int, int]) []int {
	vals := drainIterator(t, it)
	for i := 0; i < len(vals)/2; i++ {
		vals[i] = vals[i*2]
	}
	return vals[:len(vals)/2]
}

// drainIterator returns slice of key-value pairs of [K1, V1, K2, V2, ...]
func drainIterator(t *testing.T, it Iterator[int, int]) []int {
	var result []int
	for it.HasNext() {
		k, v, err := it.Next()
		require.NoError(t, err)
		result = append(result, k, v)
	}
	_, _, err := it.Next()
	require.ErrorIs(t, err, ErrEmptyIterator)
	return result
}
