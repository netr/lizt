package lizt_test

import (
	"testing"

	"git.faze.center/netr/lizt"
)

func TestCountLines(t *testing.T) {
	count, err := lizt.FileLineCount("test/1000000.txt")
	if err != nil {
		t.Error(err)
	}

	if count != 1000001 {
		t.Error("Expected 1000001 lines, got", count)
	}
}

func TestRepeatLines(t *testing.T) {
	lines, err := lizt.ReadFromFile("test/1000000.txt")
	if err != nil {
		t.Error(err)
	}

	count := 1000001
	times := 10
	rep := lizt.RepeatLines(lines, times)
	if len(rep) != count*times {
		t.Error("Expected", count*times, "lines, got", len(rep))
	}
}

// BenchmarkReadFromFile
// BenchmarkReadFromFile-8           	       7	 164656063 ns/op
// BenchmarkReadFromFilePreAlloc
// BenchmarkReadFromFilePreAlloc-8   	      24	  49227811 ns/op
// I removed the non-pre-allocated version because it was slower
// Even though the pre-allocated version required opening the file *twice*,
// It's still faster than the non-pre-allocated version.
func BenchmarkReadFromFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := lizt.ReadFromFile("test/1000000.txt")
		if err != nil {
			b.Error(err)
		}
	}
}
