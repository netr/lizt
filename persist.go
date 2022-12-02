package lizt

type PersistentIterator struct {
	Persister
	PointerIterator
}

type PersistentSeederIterator struct {
	Seeder
	PersistentIterator
}
