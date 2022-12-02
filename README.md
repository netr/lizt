# Lizt
lizt is a flexible list/file manager. you can create stream iterators or slice iterators as the base. there is a persistent storage wrapper and a seeder wrapper.

## Usage

### Slice Iterator With Seeder Wrapper
```go
package main
import "git.faze.center/netr/lizt"

func main() {
	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	seedIter := lizt.NewSliceIterator("seeds", []string{"seeder1", "seeder2"}, true)

	seed := lizt.NewSeedingIterator(
		lizt.NewSliceIterator("numbers", numbers, false),
		seedIter, // useful for sharing seeder across multiple lists  
		2,        // seed after
	)

	_, err = seed.Next(4)
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

func main() {
	seedIter := lizt.NewSliceIterator("seeds", []string{"seeder1", "seeder2"}, true)

	roundRobin := true
	stream, err := lizt.NewStreamIterator("filename.txt", roundRobin)
	if err != nil {
		panic(err)
	}

	seed := lizt.NewSeedingIterator(
		stream,
		seedIter, // useful for sharing seeder across multiple lists  
		2,        // seed after
	)

	_, err = seed.Next(4)
	if err != nil {
		panic(err)
	}
	// results in ["{line_1}", "seeder1", "{line_2}", "seeder2"]
}
```


