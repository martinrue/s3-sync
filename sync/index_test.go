package sync_test

import (
	"testing"

	"github.com/martinrue/s3-sync/sync"
)

type input struct {
	Key  string
	Tags string
}

type expectation struct {
	Key   string
	Tags  string
	IsNew bool
}

type test struct {
	Name     string
	Left     []input
	Right    []input
	Expected []expectation
}

func TestIndexDiff(t *testing.T) {
	tests := []test{
		{
			Name:     "Empty indexes, empty diff",
			Left:     []input{},
			Right:    []input{},
			Expected: []expectation{},
		},
		{
			Name: "Equal indexes, empty diff",
			Left: []input{
				{"key-1", "one unu"},
				{"key-2", "two du"},
			},
			Right: []input{
				{"key-2", "two du"},
				{"key-1", "one unu"},
			},
			Expected: []expectation{},
		},
		{
			Name: "Unequal indexes, two new keys",
			Left: []input{
				{"key-1", "one unu"},
				{"key-2", "two du"},
			},
			Right: []input{
				{"key-1", "one unu"},
				{"key-3", "three tri"},
				{"key-4", "four kvar"},
			},
			Expected: []expectation{
				{"key-3", "three tri", true},
				{"key-4", "four kvar", true},
			},
		},
		{
			Name: "Unequal indexes, one new key, one updated tag",
			Left: []input{
				{"key-1", "one unu"},
				{"key-2", "two du"},
			},
			Right: []input{
				{"key-2", "two du new-tag"},
				{"key-3", "three tri"},
			},
			Expected: []expectation{
				{"key-2", "two du new-tag", false},
				{"key-3", "three tri", true},
			},
		},
		{
			Name: "Unequal indexes, no new keys",
			Left: []input{
				{"key-1", "one unu"},
				{"key-2", "two du"},
				{"key-3", "three tri"},
				{"key-4", "four kvar"},
			},
			Right: []input{
				{"key-1", "one unu"},
				{"key-3", "three tri"},
			},
			Expected: []expectation{},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			left := &sync.Index{}

			for _, data := range test.Left {
				left.Add(data.Key, data.Tags, "", false)
			}

			right := &sync.Index{}

			for _, data := range test.Right {
				right.Add(data.Key, data.Tags, "", false)
			}

			diff := left.Diff(right)

			if len(diff.Objects) != len(test.Expected) {
				t.Fatalf("diff len (%d) not equal to expected len (%d)", len(diff.Objects), len(test.Expected))
			}

			for i, expected := range test.Expected {
				actual := diff.Objects[i]

				if actual.Key != expected.Key {
					t.Fatalf("expected diff object %d to be %v, but got %v", i, expected.Key, actual.Key)
				}

				if actual.Tags != expected.Tags {
					t.Fatalf("expected diff object %d to have tags %v, but got %v", i, expected.Tags, actual.Tags)
				}

				if actual.IsNew != expected.IsNew {
					t.Fatalf("expected diff object %d to have IsNew %v, but got %v", i, expected.IsNew, actual.IsNew)
				}
			}
		})
	}
}

func TestIndexSaveJSONWithEmptyIndex(t *testing.T) {
	index := &sync.Index{}

	data, err := index.SaveJSON()
	if err != nil {
		t.Fatal(err)
	}

	expected := "null"
	actual := string(data)

	if actual != expected {
		t.Fatalf("expected data to be (%v), got (%v)", expected, actual)
	}
}

func TestIndexSaveJSONWithNonEmptyIndex(t *testing.T) {
	index := &sync.Index{
		Objects: []*sync.Object{
			&sync.Object{Key: "key-1", Tags: "one unu"},
			&sync.Object{Key: "key-2", Tags: "two du"},
		},
	}

	data, err := index.SaveJSON()
	if err != nil {
		t.Fatal(err)
	}

	expected := `[{"key":"key-1","tags":"one unu"},{"key":"key-2","tags":"two du"}]`
	actual := string(data)

	if actual != expected {
		t.Fatalf("expected data to be (%v), got (%v)", expected, actual)
	}
}

func TestIndexLoadJSON(t *testing.T) {
	index := &sync.Index{}

	data := []byte(`[{"key":"key-1","tags":"one"},{"key":"key-2","tags":"two"}]`)

	if err := index.LoadJSON(data); err != nil {
		t.Fatal(err)
	}

	if len(index.Objects) != 2 {
		t.Fatalf("expected index to have 2 objects, but has %v", len(index.Objects))
	}

	first, second := index.Objects[0], index.Objects[1]

	if first.Key != "key-1" || first.Tags != "one" {
		t.Fatalf("unexpected first object %+v", first)
	}

	if second.Key != "key-2" || second.Tags != "two" {
		t.Fatalf("unexpected second object %+v", first)
	}
}
