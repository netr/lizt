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

	if count != 999998 {
		t.Error("Expected 999998 lines, got", count)
	}
}

func TestRepeatLines(t *testing.T) {
	lines, err := lizt.ReadFromFile("test/1000000.txt")
	if err != nil {
		t.Error(err)
	}

	count := 999998
	times := 10
	rep := lizt.RepeatLines(lines, times)
	if len(rep) != count*times {
		t.Error("Expected", count*times, "lines, got", count)
	}
}
