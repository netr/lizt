# Lizt
lizt is a flexible list/file manager. you can create stream iterators or slice iterators as the base. there is a persistent storage wrapper and a seeder wrapper.

## Read first
Based on the benchmarks, it seems that the SliceIterator is faster in every case. Which is pretty much true, once the slice is in memory.

The major differences between the two iterators are:
- Slice requires holding the entire list in memory and is therefore limited by the amount of memory available. Stream does not have this limitation.
- Slice is faster when the list is small. Stream is faster when the list is large.
- Slices will break eventually. For really large lists (10M+), streams are always the better choice.
- Streams should **not be used at all** if the list is small. The overhead of re-creating a stream, when round-robin is being used, is not worth it. Instead, `lizt.ReadFromFile(path)` into a SliceIterator.

## Usage

### Builder Helper Examples

#### File Stream Iterator
```go
stream, _ := lizt.NewBuilder().Stream("test/50000000.txt").Build() // round-robin = false

fmt.Println(stream.Next(5)) // "a", "b", "c", "e", "f"

stream, _ := lizt.NewBuilder().StreamRR("test/50000000.txt").Build() // round-robin = true

fmt.Println(stream.Next(5)) // "a", "b", "c", "e", "f"
```

#### Slice Iterator
```go
// creates a random string for it's name for ease of use

slice, _ := lizt.B().Slice([]string{"a", "b", "c", "d", "e"}).Build() // round-robin = false

fmt.Println(slice.Next(3)) // "a", "b", "c"

slice, _ := lizt.B().SliceRR([]string{"a", "b", "c", "d", "e"}).Build() // round-robin = true

fmt.Println(slice.Next(3)) // "a", "b", "c"
```

#### Slice Iterator With Name (Required if you're using a Manager)
```go
slice, _ := lizt.B().SliceNamed("name", []string{"a", "b", "c", "d", "e"}).Build() // round-robin = false

fmt.Println(slice.Next(3)) // "a", "b", "c"

slice, _ := lizt.B().SliceNamedRR("name", []string{"a", "b", "c", "d", "e"}).Build() // round-robin = true

fmt.Println(slice.Next(3)) // "a", "b", "c"
```

#### Seeding File Stream Iterator with Seed Iterator as a Slice
```go
plantEvery := 2

seedStream, _ = lizt.B().
            Stream("test/10000000.txt").
            BuildWithSeeds(plantEvery, []string{"seed1", "seed2"}) // round-robin = false

seedStream, _ = lizt.B().
            StreamRR("test/50000000.txt").
            BuildWithSeeds(plantEvery, []string{"seed1", "seed2"}) // round-robin = true

fmt.Println(slice.Next(4)) // "seed1", "a", "seed2", "b"
```

#### Seeding Slice Iterator with Seed Iterator as a File Stream
```go
plantEvery := 2

seedStream, _ = lizt.B().
            Stream("test/10000000.txt").
            BuildWithSeeds(plantEvery, "data/seeds.txt") // round-robin = false

seedStream, _ = lizt.B().
            StreamRR("test/50000000.txt").
            BuildWithSeeds(plantEvery, "data/seeds.txt") // round-robin = true

fmt.Println(slice.Next(4))// "seed1", "a", "seed2", "b"
	
```

#### File Stream Iterator with Persistence ( See `persist_test.go` for more examples )
```go
mem := NewInMemoryPersister()
stream, _ := lizt.B().StreamRR("test/50000000.txt").PersistTo(mem).Build() // round-robin = false

fmt.Println(stream.Next(5)) // "a", "b", "c", "e", "f"
```

#### Slice Iterator with Persistence
```go
mem := NewInMemoryPersister()

stream, _ := lizt.B().Slice([]{"test","this","here"}).PersistTo(mem).Build() // round-robin = false

fmt.Println(stream.Next(3)) // "test", "this", "here"
// mem["10"] = 3
```

#### Seeding File Stream Iterator with Persistence
```go
mem := NewInMemoryPersister()

stream, _ := lizt.B().
        StreamRR("test/50000000.txt").
        PersistTo(mem).
        BuildWithSeeds(2, []string{"seed1", "seed2"}) // round-robin = false
		
fmt.Println(stream.Next(5)) // "seed1", "a", "seed2", "b", "seed1"
// Persister Value => mem["10"] = 2
```

## Using the Manager

```go
package main

import "git.faze.center/netr/lizt"

type IterKey string

const (
	IterKeyFiftyMillion IterKey = "50000000"
	IterKeyNumbers      IterKey = "numbers"
)

func main() {
	mem := NewInMemoryPersister()

	// B() => NewBuilder(), StreamRR() => Stream with Round Robin, PersistTo() => Persist to Persister, BuildWithSeeds() => Build the Iterator with Seeding
	fiftyStream, err := lizt.B().StreamRR("test/50000000.txt").PersistTo(mem).BuildWithSeeds(2, []string{"seeder1", "seeder2"}) // round-robin = false
	if err != nil {
		panic(err)
	}

	// B() => NewBuilder(), StreamRR() => Stream with Round Robin, PersistTo() => Persist to Persister, Build() => Build the Iterator
	numbers := []string{"1", "2", "3", "4", "5"}
	numbersSlice, err := lizt.B().SliceNamedRR(string(IterKeyNumbers), numbers).PersistTo(mem).Build() // round-robin = false
	if err != nil {
		panic(err)
	}

	// initialize and add seeder to manager
	mgr := lizt.NewManager().AddIters(fiftyStream, numbersSlice)

	_, err = mgr.MustGet(string(IterKeyFiftyMillion)).Next(4)
	if err != nil {
		panic(err)
	}
	// results in ["seeder1", "1", "seeder2", "2"]

	_, err = mgr.MustGet(string(IterKeyNumbers)).Next(4)
	if err != nil {
		panic(err)
	}
	// results in ["1", "2", "3", "4"]
}
```
