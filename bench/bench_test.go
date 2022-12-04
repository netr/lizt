package bench_test

import (
	"git.faze.center/netr/lizt"
	"testing"
	"time"
)

var (
	filenameOneMillion   = "../test/1000000.txt"
	filenameTenMillion   = "10000000.txt"
	filenameFiftyMillion = "50000000.txt"
	filenameTen          = "../test/10.txt"
	useTimeReporting     = false
)

func Test_CreateLargeSeedData(t *testing.T) {
	disabled := true
	if !disabled {
		lines, err := lizt.ReadFromFile(filenameOneMillion)
		if err != nil {
			t.Errorf("ReadFromFile() error = %v", err)
		}
		err = lizt.WriteToFile(lizt.RepeatLines(lines, 10), filenameTenMillion)
		if err != nil {
			t.Errorf("WriteToFile() error = %v", err)
		}
	}
}

func Test_CreateLargestSeedData(t *testing.T) {
	disabled := true
	if !disabled {
		lines, err := lizt.ReadFromFile(filenameOneMillion)
		if err != nil {
			t.Errorf("ReadFromFile() error = %v", err)
		}
		err = lizt.WriteToFile(lizt.RepeatLines(lines, 50), filenameFiftyMillion)
		if err != nil {
			t.Errorf("WriteToFile() error = %v", err)
		}
	}
}

func BenchmarkStreamIterator_Next_10(b *testing.B) {
	start := time.Now()

	mgr := lizt.NewManager()
	fs, err := lizt.B().StreamRR(filenameTen).Build()
	if err != nil {
		b.Errorf("NewStreamIterator() error = %v", err)
	}
	err = mgr.AddIter(fs)
	if err != nil {
		b.Errorf("AddIter() error = %v", err)
	}

	since := time.Since(start)
	b.Logf("BenchmarkStreamIterator_Next_10: Setup took: %v", since)

	for i := 0; i < b.N; i++ {
		_, err = mgr.MustGet("10").Next(10)
		if err != nil {
			b.Errorf("StreamIterator.Next() error = %v", err)
		}
	}
}

func BenchmarkSliceIterator_Next_10(b *testing.B) {
	start := time.Now()

	mgr := lizt.NewManager()
	slice, err := lizt.ReadFromFile(filenameTen)
	if err != nil {
		b.Errorf("ReadFromFile() error = %v", err)
	}
	fs, err := lizt.B().SliceNamedRR("10", slice).Build()
	if err != nil {
		b.Errorf("NewStreamIterator() error = %v", err)
	}
	err = mgr.AddIter(fs)
	if err != nil {
		b.Errorf("AddIter() error = %v", err)
	}

	since := time.Since(start)
	b.Logf("BenchmarkSliceIterator_Next_10: Setup took: %v", since)

	for i := 0; i < b.N; i++ {
		_, err = mgr.MustGet("10").Next(10)
		if err != nil {
			b.Errorf("StreamIterator.Next() error = %v", err)
		}
	}
}

func BenchmarkStreamIterator_Next_1000000(b *testing.B) {
	start := time.Now()

	mgr := lizt.NewManager()
	fs, err := lizt.B().StreamRR(filenameOneMillion).Build()
	if err != nil {
		b.Errorf("NewStreamIterator() error = %v", err)
	}
	err = mgr.AddIter(fs)
	if err != nil {
		b.Errorf("AddIter() error = %v", err)
	}

	since := time.Since(start)
	b.Logf("BenchmarkStreamIterator_Next_1000000: Setup took: %v", since)

	for i := 0; i < b.N; i++ {
		_, err = mgr.MustGet("1000000").Next(10)
		if err != nil {
			b.Errorf("StreamIterator.Next() error = %v", err)
		}
	}
}

func BenchmarkSliceIterator_Next_1000000(b *testing.B) {
	start := time.Now()

	mgr := lizt.NewManager()
	slice, err := lizt.ReadFromFile(filenameOneMillion)
	if err != nil {
		b.Errorf("ReadFromFile() error = %v", err)
	}
	fs, err := lizt.B().SliceNamedRR("1000000", slice).Build()
	if err != nil {
		b.Errorf("NewStreamIterator() error = %v", err)
	}
	err = mgr.AddIter(fs)
	if err != nil {
		b.Errorf("AddIter() error = %v", err)
	}

	since := time.Since(start)
	b.Logf("BenchmarkSliceIterator_Next_1000000: Setup took: %v", since)

	for i := 0; i < b.N; i++ {
		_, err = mgr.MustGet("1000000").Next(10)
		if err != nil {
			b.Errorf("StreamIterator.Next() error = %v", err)
		}
	}
}

func BenchmarkStreamIterator_Next_10000000(b *testing.B) {
	start := time.Now()

	mgr := lizt.NewManager()
	fs, err := lizt.B().StreamRR(filenameTenMillion).Build()
	if err != nil {
		b.Errorf("NewStreamIterator() error = %v", err)
	}
	err = mgr.AddIter(fs)
	if err != nil {
		b.Errorf("AddIter() error = %v", err)
	}

	since := time.Since(start)
	b.Logf("BenchmarkStreamIterator_Next_10000000: Setup took %v", since)

	for i := 0; i < b.N; i++ {
		_, err = mgr.MustGet("10000000").Next(10)
		if err != nil {
			b.Errorf("StreamIterator.Next() error = %v", err)
		}
	}
}

func BenchmarkSliceIterator_Next_10000000(b *testing.B) {
	start := time.Now()

	mgr := lizt.NewManager()
	slice, err := lizt.ReadFromFile(filenameTenMillion)
	if err != nil {
		b.Errorf("ReadFromFile() error = %v", err)
	}
	fs, err := lizt.B().SliceNamedRR("10000000", slice).Build()
	if err != nil {
		b.Errorf("NewStreamIterator() error = %v", err)
	}
	err = mgr.AddIter(fs)
	if err != nil {
		b.Errorf("AddIter() error = %v", err)
	}

	since := time.Since(start)
	b.Logf("BenchmarkSliceIterator_Next_10000000: Setup took: %v", since)

	for i := 0; i < b.N; i++ {
		_, err = mgr.MustGet("10000000").Next(10)
		if err != nil {
			b.Errorf("StreamIterator.Next() error = %v", err)
		}
	}
}

func BenchmarkStreamIterator_Next_50000000(b *testing.B) {
	start := time.Now()

	mgr := lizt.NewManager()
	fs, err := lizt.B().StreamRR(filenameFiftyMillion).Build()
	if err != nil {
		b.Errorf("NewStreamIterator() error = %v", err)
	}
	err = mgr.AddIter(fs)
	if err != nil {
		b.Errorf("AddIter() error = %v", err)
	}

	since := time.Since(start)
	b.Logf("BenchmarkStreamIterator_Next_50000000: Setup took: %v", since)

	for i := 0; i < b.N; i++ {
		_, err = mgr.MustGet("50000000").Next(10)
		if err != nil {
			b.Errorf("StreamIterator.Next() error = %v", err)
		}
	}
}

func BenchmarkSliceIterator_Next_50000000(b *testing.B) {
	start := time.Now()

	mgr := lizt.NewManager()
	slice, err := lizt.ReadFromFile(filenameFiftyMillion)
	if err != nil {
		b.Errorf("ReadFromFile() error = %v", err)
	}
	fs, err := lizt.B().SliceNamedRR("50000000", slice).Build()
	if err != nil {
		b.Errorf("NewStreamIterator() error = %v", err)
	}
	err = mgr.AddIter(fs)
	if err != nil {
		b.Errorf("AddIter() error = %v", err)
	}

	since := time.Since(start)
	b.Logf("BenchmarkSliceIterator_Next_50000000: Setup took: %v", since)

	for i := 0; i < b.N; i++ {
		_, err = mgr.MustGet("50000000").Next(10)
		if err != nil {
			b.Errorf("StreamIterator.Next() error = %v", err)
		}
	}
}
