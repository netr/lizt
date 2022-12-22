package lizt_test

import (
	"testing"

	"git.faze.center/netr/lizt"
)

func TestBlacklister_Next(t *testing.T) {
	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	blacklist := lizt.BlacklistMap{"2": {}, "4": {}, "6": {}, "8": {}, "10": {}}
	blm := lizt.NewBlacklistManager(blacklist)

	blkIter, _ := lizt.B().Slice(numbers).Blacklist(blm).Build()

	next, err := blkIter.Next(5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := map[int]string{
		0: "1",
		1: "3",
		2: "5",
		3: "7",
		4: "9",
	}

	for k, v := range expected {
		if next[k] != v {
			t.Errorf("Expected %s, got %s", v, next[k])
		}
	}
}

func TestBlacklister_Next_ShouldNotReturnZeroEntriesIfItRemovesAllOfThem(t *testing.T) {
	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	blacklist := map[string]struct{}{"1": {}, "2": {}, "3": {}, "4": {}, "5": {}}
	blm := lizt.NewBlacklistManager(blacklist)

	blkIter, _ := lizt.B().Slice(numbers).Blacklist(blm).Build()

	next, err := blkIter.Next(5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := map[int]string{
		0: "6",
		1: "7",
		2: "8",
		3: "9",
		4: "10",
	}

	for k, v := range expected {
		if next[k] != v {
			t.Errorf("Expected %s, got %s", v, next[k])
		}
	}
}

func TestScrubFileWithBlacklist(t *testing.T) {
	blkMap := lizt.BlacklistMap{
		"b": {}, "d": {}, "f": {}, "h": {}, "j": {},
	}

	// blkMap, err := lizt.FileTomap("test/blacklist.txt", blkMap)

	n, err := lizt.ScrubFileWithBlacklist(blkMap, "test/10.txt", "test/10.txt.scrubbed")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if n != 5 {
		t.Fatalf("Expected 5 lines scrubbed, got %d", n)
	}

	scrubbed, err := lizt.ReadFromFile("test/10.txt.scrubbed")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := []string{"a", "c", "e", "g", "i"}
	for k, v := range expected {
		if scrubbed[k] != v {
			t.Errorf("Expected %s, got %s", v, scrubbed[k])
		}
	}

	err = lizt.DeleteFile("test/10.txt.scrubbed")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
