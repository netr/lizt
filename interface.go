package lizt

// Iterator is an interface for iterating over a list of lines
type Iterator interface {
	Name() string
	Len() int
	Next(count int) ([]string, error)
}

// PointerIterator is an iterator that has a pointer
type PointerIterator interface {
	Iterator
	Pointer() uint64
	Inc()
	ResetPointer()
}

// Seeder is an interface for seeding a pointer iterator
type Seeder interface {
	Seeds() []string
	PlantEvery() int
	Planted() int64
}

// Persister adds persistent storage to an iterator
type Persister interface {
	Set(key, value string) error
	Get(key string) (string, error)
}
