## Benchmark Results

```go
goos: linux
goarch: amd64
pkg: git.faze.center/netr/lizt/bench
cpu: AMD Ryzen 9 3950X 16-Core Processor
BenchmarkStreamIterator_Next_10
2022/12/04 09:38:24 BenchmarkStreamIterator_Next_10: Setup took 76.232µs
BenchmarkStreamIterator_Next_10-32          	  150440	      8028 ns/op
BenchmarkSliceIterator_Next_10
2022/12/04 09:38:26 BenchmarkSliceIterator_Next_10: Setup took 26.891µs
BenchmarkSliceIterator_Next_10-32           	 3178616	       372.5 ns/op
BenchmarkStreamIterator_Next_1000000
2022/12/04 09:38:27 BenchmarkStreamIterator_Next_1000000: Setup took 13.208939ms
BenchmarkStreamIterator_Next_1000000-32     	 1225536	       966.3 ns/op
BenchmarkSliceIterator_Next_1000000
2022/12/04 09:38:29 BenchmarkSliceIterator_Next_1000000: Setup took 59.131062ms
BenchmarkSliceIterator_Next_1000000-32      	 2552266	       433.4 ns/op
BenchmarkStreamIterator_Next_10000000
2022/12/04 09:38:31 BenchmarkStreamIterator_Next_10000000: Setup took 127.127557ms
BenchmarkStreamIterator_Next_10000000-32    	 1183899	       981.1 ns/op
BenchmarkSliceIterator_Next_10000000
2022/12/04 09:38:35 BenchmarkSliceIterator_Next_10000000: Setup took 511.000518ms
BenchmarkSliceIterator_Next_10000000-32     	 1339508	       784.8 ns/op
BenchmarkStreamIterator_Next_50000000
2022/12/04 09:38:48 BenchmarkStreamIterator_Next_50000000: Setup took 635.567025ms
BenchmarkStreamIterator_Next_50000000-32    	  479552	      2212 ns/op
BenchmarkSliceIterator_Next_50000000
2022/12/04 09:39:07 BenchmarkSliceIterator_Next_50000000: Setup took 2.501227454s
BenchmarkSliceIterator_Next_50000000-32     	       1	2501257505 ns/op
PASS
```
----

*Note:* To generate seed data for tests uncomment `t.SkipNow()` in the two `Test_CreateLargeSeedData`  and `Test_CreateLargestSeedData` functions

----

See `README.md` for breakdown.