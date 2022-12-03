# Lizt
lizt is a flexible list/file manager. you can create stream iterators or slice iterators as the base. there is a persistent storage wrapper and a seeder wrapper.

## Usage

### Slice Iterator With Seeder Wrapper
```go
package main
import "git.faze.center/netr/lizt"

type IterKey string
const (
	IterKeyNumber IterKey = "numbers"   
	IterKeySeeds IterKey = "seeds"
)

func main() {
    numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
    numbersIter := lizt.NewSliceIterator(IterKeyNumber, numbers, false)
    seedIter := lizt.NewSliceIterator(IterKeySeeds, []string{"seeder1", "seeder2"}, true)

    seed := lizt.NewSeedingIterator(
        lizt.SeedingIteratorConfig{
            PointerIterator: numbersIter,
            Seeds:           seedIter,
            SeedEvery:       2,
        },
    )
    
    // initialize and add seeder to manager
    mgr := lizt.NewManager()
    err = mgr.AddIter(seed)
    if err != nil {
        panic(err)
    }
    
    _, err = mgr.Get(IterKeyNumber).Next(4)
    if err != nil {
        panic(err)
    }
    // results in ["1", "seeder1", "2", "seeder2"]
}
```


### Stream Iterator With Seeder Wrapper
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
            PointerIterator: streamIter,
            Seeds:           seedIter,
            SeedEvery:       2,
        },
    )
    
    // initialize and add seeder to manager
    mgr := lizt.NewManager()
    err = mgr.AddIter(seed)
    if err != nil {
        panic(err)
    }
    
    _, err = mgr.Get(IterKeyExample).Next(4)
    if err != nil {
        panic(err)
    }
    // results in ["{line_1}", "seeder1", "{line_2}", "seeder2"]
}
```


