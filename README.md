# Lizt
lizt is a flexible list/file manager. you can create stream iterators or slice iterators as the base. there is a persistent storage wrapper and a seeder wrapper.

## Usage

## Builder Helper Examples

#### File Stream Iterator
```go
stream, _ := lizt.NewBuilder().Stream("test/10.txt").Build() // round-robin = false
fmt.Println(stream.Next(5))
// "a", "b", "c", "e", "f"

stream, _ := lizt.NewBuilder().StreamRR("test/10.txt").Build() // round-robin = true
fmt.Println(stream.Next(5))
// "a", "b", "c", "e", "f"
```

#### Slice Iterator
```go
// creates a random string for it's name for ease of use

slice, _ := lizt.B().Slice([]string{"a", "b", "c", "d", "e"}).Build() // round-robin = false
fmt.Println(slice.Next(3))

slice, _ := lizt.B().SliceRR([]string{"a", "b", "c", "d", "e"}).Build() // round-robin = true
fmt.Println(slice.Next(3))
// "a", "b", "c"
```

#### Slice Iterator With Name (Required if you're using a Manager)
```go
slice, _ := lizt.B().SliceNamed("name", []string{"a", "b", "c", "d", "e"}).Build() // round-robin = false
fmt.Println(slice.Next(3))
// "a", "b", "c"

slice, _ := lizt.B().SliceNamedRR("name", []string{"a", "b", "c", "d", "e"}).Build() // round-robin = true
fmt.Println(slice.Next(3))
// "a", "b", "c"
```

#### Seeding File Stream Iterator with Seed Iterator as a Slice
```go
plantEvery := 2
seedStream, _ = lizt.B().
            Stream("test/10.txt").
            BuildWithSeeds(plantEvery, []string{"seed1", "seed2"}) // round-robin = false

seedStream, _ = lizt.B().
            StreamRR("test/10.txt").
            BuildWithSeeds(plantEvery, []string{"seed1", "seed2"}) // round-robin = true

fmt.Println(slice.Next(4))
// "seed1", "a", "seed2", "b"
```

#### Seeding Slice Iterator with Seed Iterator as a File Stream
```go
plantEvery := 2
seedStream, _ = lizt.B().
            Stream("test/10.txt").
            BuildWithSeeds(plantEvery, "data/seeds.txt") // round-robin = false

seedStream, _ = lizt.B().
            StreamRR("test/10.txt").
            BuildWithSeeds(plantEvery, "data/seeds.txt") // round-robin = true

fmt.Println(slice.Next(4))
// "seed1", "a", "seed2", "b"
	
```

#### File Stream Iterator with Persistence ( See `persist_test.go` for more examples )
```go
mem := NewInMemoryPersister()
stream, _ := lizt.NewBuilder().Stream("test/10.txt").PersistTo(mem).Build() // round-robin = false
fmt.Println(stream.Next(5))
// "a", "b", "c", "e", "f"

// See `persist_test.go` for more examples
mem := NewInMemoryPersister()
stream, _ := lizt.NewBuilder().StreamRR("test/10.txt").PersistTo(mem).Build() // round-robin = true
fmt.Println(stream.Next(5))
// "a", "b", "c", "e", "f"
```

#### Slice Iterator with Persistence ( See `persist_test.go` for more examples )
```go
mem := NewInMemoryPersister()
stream, _ := lizt.NewBuilder().Slice([]{"test","this","here"}).PersistTo(mem).Build() // round-robin = false
fmt.Println(stream.Next(3))
// "test", "this", "here"
// mem["10"] = 3

// See `persist_test.go` for more examples
mem := NewInMemoryPersister()
stream, _ := lizt.NewBuilder().SliceRR([]{"test","this","here"}).PersistTo(mem).Build() // round-robin = true
fmt.Println(stream.Next(6))
// "test", "this", "here", "test", "this", "here"
// mem["10"] = 3
```

## Slice Iterator With SeedingIterator Wrapper
```go
package main
import "git.faze.center/netr/lizt"

type IterKey string
const (
	IterKeyNumbers IterKey = "numbers"   
	IterKeySeeds IterKey = "seeds"
)

func main() {
    numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
    numbersIter := lizt.NewSliceIterator(IterKeyNumbers, numbers, false)
    seedIter := lizt.NewSliceIterator(IterKeySeeds, []string{"seeder1", "seeder2"}, true)

    seed := lizt.NewSeedingIterator(
        lizt.SeedingIteratorConfig{
            PointerIter:   numbersIter,
            SeedIter:      seedIter,
            PlantEvery:    2,
        },
    )
    
    // initialize and add seeder to manager
    mgr := lizt.NewManager()
    err = mgr.AddIter(seed)
    if err != nil {
        panic(err)
    }
    
    _, err = mgr.MustGet(IterKeyNumbers).Next(4)
    if err != nil {
        panic(err)
    }
    // results in ["seeder1", "1", "seeder2", "2"]
}
```

## File Stream Iterator With SeedingIterator Wrapper
```go
package main
import "git.faze.center/netr/lizt"

type IterKey string
const (
    IterKeyExample IterKey = "example"
    IterKeySeeds IterKey = "seeds"
)

func main() {
    roundRobin := true
    seedIter := lizt.NewSliceIterator(IterKeySeeds, []string{"seeder1", "seeder2"}, roundRobin)
    streamIter, err := lizt.NewStreamIterator("example.txt", roundRobin)
	
    if err != nil {
        panic(err)
    }
    
    seed := lizt.NewSeedingIterator(
        lizt.SeedingIteratorConfig{
	        PointerIter:   streamIter,
	        SeedIter:      seedIter,
	        PlantEvery:    2,
        },
    )
    
    // initialize and add seeder to manager
    mgr := lizt.NewManager()
    err = mgr.AddIter(seed)
    if err != nil {
        panic(err)
    }
    
    _, err = mgr.MustGet(IterKeyExample).Next(4)
    if err != nil {
        panic(err)
    }
    // results in ["seeder1", "{line_1}", "seeder2", "{line_2}"]
}
```

## Persistent Wrapper
See `TestPersistentIterator_Next` and `TestNewPersistentIterator_UsingSeedingIterator_Next` in `persist_test.go` for examples.
