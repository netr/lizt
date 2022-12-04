# Lizt
lizt is a flexible list/file manager. you can create stream iterators or slice iterators as the base. there is a persistent storage wrapper and a seeder wrapper.

## Usage

## Builder Helper Examples

#### Stream Iterator
```go
stream, _ := lizt.Builder().Stream("test/10.txt", roundRobin).Build()
fmt.Println(stream.Next(5))
// "a", "b", "c", "e", "f"
```

#### Slice Iterator
```go
slice, _ := lizt.Builder().Slice("name", []string{"a", "b", "c", "d", "e"}, roundRobin).Build()
fmt.Println(slice.Next(3))
// "a", "b", "c"
```

#### Seeding Stream Iterator with Seed Slice Iterator
```go
plantEvery := 2
seedStream, _ = lizt.Builder().
            Stream("test/10.txt", roundRobin).
            BuildWithSeeds(plantEvery, []string{"seed1", "seed2"})

fmt.Println(slice.Next(4))
// "seed1", "a", "seed2", "b"
```

#### Seeding Slice Iterator with Seed Stream Iterator
```go
plantEvery := 2
seedStream, _ = lizt.Builder().
            Stream("test/10.txt", roundRobin).
            BuildWithSeeds(plantEvery, "data/seeds.txt")

fmt.Println(slice.Next(4))
// "seed1", "a", "seed2", "b"
	
```

## Slice Iterator With Seeder Wrapper
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

## Stream Iterator With Seeder Wrapper
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
