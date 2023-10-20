package set

type Set[T comparable] struct {
	data             map[T]bool
	optFlattenedKeys []T
}

func (s *Set[T]) Add(item T) {
	s.data[item] = true
}

func (s *Set[T]) Remove(item T) {
	delete(s.data, item)
}

func (s *Set[T]) Has(item T) bool {
	return s.data[item]
}

func (s *Set[T]) Size() int {
	return len(s.data)
}

func (s *Set[T]) Flatten() []T {
	if s.optFlattenedKeys != nil {
		return s.optFlattenedKeys
	}
	keys := make([]T, len(s.data))
	keysIdx := 0
	for key, _ := range s.data {
		keys[keysIdx] = key
		keysIdx++
	}
	s.optFlattenedKeys = keys
	return keys
}

func EmptySet[T comparable]() *Set[T] {
	return &Set[T]{
		data: make(map[T]bool),
	}
}
