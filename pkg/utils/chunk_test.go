package utils

import (
	"reflect"
	"testing"
)

func TestChunkStringSlice(t *testing.T) {
	in := []string{"a", "b", "c", "d", "e"}
	expected := [][]string{
		{
			"a",
			"b",
		},
		{
			"c",
			"d",
		},
		{
			"e",
		},
	}

	out, err := ChunkStringSlice(in, 2)
	if err != nil {
		t.Fatal("must not return error")
	}

	if !reflect.DeepEqual(out, expected) {
		t.Fatal("must be equal")
	}
}

func TestChunkStringSlice2(t *testing.T) {
	in := []string{"a", "b", "c", "d", "e"}
	expected := [][]string{
		{
			"a",
			"b",
			"c",
			"d",
			"e",
		},
	}

	out, err := ChunkStringSlice(in, 6)
	if err != nil {
		t.Fatal("must not return error")
	}

	if !reflect.DeepEqual(out, expected) {
		t.Fatal("must be equal")
	}
}

func TestChunkStringSlice3(t *testing.T) {
	in := []string{"a", "b", "c", "d", "e"}
	expected := [][]string{
		{
			"a",
			"b",
			"c",
			"d",
			"e",
		},
	}

	out, err := ChunkStringSlice(in, 6)

	if err != nil {
		t.Fatal("must not return error")
	}

	if !reflect.DeepEqual(out, expected) {
		t.Fatal("must be equal")
	}
}

func TestChunkStringSliceError(t *testing.T) {
	in := []string{"a", "b", "c", "d", "e"}
	_, err := ChunkStringSlice(in, 0)

	if err == nil {
		t.Fatal("must return error")
	}
}
