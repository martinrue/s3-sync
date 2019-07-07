package sync_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/martinrue/s3-sync/sync"
)

type Object struct {
	ContentType string
	Data        []byte
}

type MockObjectStore struct {
	Objects map[string]*Object
	Gets    int
	Puts    int
}

func (m *MockObjectStore) Get(key string) ([]byte, error) {
	m.Gets++

	object, ok := m.Objects[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %v", key)
	}

	return object.Data, nil

}

func (m *MockObjectStore) Put(key string, contentType string, data io.ReadSeeker) error {
	bytes, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	m.Puts++
	m.Objects[key] = &Object{contentType, bytes}

	return nil
}

type Assertion struct {
	Key  string
	Tags string
}

func checkAssertions(assertions []Assertion, index *sync.Index, t *testing.T) {
	find := func(key string) *sync.Object {
		for _, object := range index.Objects {
			if object.Key == key {
				return object
			}
		}

		return nil
	}

	for i, assertion := range assertions {
		object := find(assertion.Key)

		if object == nil {
			t.Fatalf("expected object at %v, but didn't find one", i)
		}

		if object.Tags != assertion.Tags {
			t.Fatalf("expected object %v to have tags: %v, but got: %v", i, assertion.Tags, object.Tags)
		}
	}
}

func setup() (*MockObjectStore, *sync.Syncer) {
	store := &MockObjectStore{
		Objects: make(map[string]*Object, 0),
	}

	syncer := &sync.Syncer{
		Store: store,
		Log:   func(message string, args ...interface{}) {},
	}

	return store, syncer
}

func TestRunNoMatchedFiles(t *testing.T) {
	_, syncer := setup()
	_, err := syncer.Run("test_data/sync_1", "none")

	expected := "no files matching (none) found in: test_data/sync_1"

	if err.Error() != expected {
		t.Fatalf("expected err: %v, but got: %v", expected, err)
	}
}

func TestRunNoExistingIndex(t *testing.T) {
	store, syncer := setup()

	json, err := syncer.Run("test_data/sync_1", "txt")
	if err != nil {
		t.Fatalf("expected err to be nil, but got: %v", err)
	}

	if store.Gets != 1 {
		t.Fatalf("expected gets to be 1, but got: %v", store.Gets)
	}

	if store.Puts != 4 {
		t.Fatalf("expected puts to be 4, but got: %v", store.Puts)
	}

	index := &sync.Index{}
	if err := index.LoadJSON([]byte(json)); err != nil {
		t.Fatalf("expected valid json result, but got err: %v", err)
	}

	assertions := []Assertion{
		{"2fa6e1156e6c6b2c808e9ef3c63f9c49ce809579200f2a40487949efd4febca4", "tag-1"},
		{"99d5293841b29576ce680d53742c480f892e83ae441919b6f4f1f590879a5a24", "tag-1 tag-2"},
		{"23b2ff797a148f61f24c08cecb8182af89ac8a761665611e2e152eda75449753", "tag-1 tag-2 tag-3"},
	}

	checkAssertions(assertions, index, t)
}

func TestRunWithExistingIndex(t *testing.T) {
	store, syncer := setup()

	_, err := syncer.Run("test_data/sync_1", "txt")
	if err != nil {
		t.Fatalf("expected first sync run err to be nil, but got: %v", err)
	}

	json, err := syncer.Run("test_data/sync_2", "txt")
	if err != nil {
		t.Fatalf("expected second sync run err to be nil, but got: %v", err)
	}

	if store.Gets != 2 {
		t.Fatalf("expected gets to be 1, but got: %v", store.Gets)
	}

	if store.Puts != 6 {
		t.Fatalf("expected puts to be 4, but got: %v", store.Puts)
	}

	index := &sync.Index{}
	if err := index.LoadJSON([]byte(json)); err != nil {
		t.Fatalf("expected valid json result, but got err: %v", err)
	}

	assertions := []Assertion{
		{"2fa6e1156e6c6b2c808e9ef3c63f9c49ce809579200f2a40487949efd4febca4", "tag-1"},
		{"37e4b113f412dc5b776b4acff3339ed4505314c8e060faecf8641e3f18f882dd", "tag-1 tag-2"},
		{"23b2ff797a148f61f24c08cecb8182af89ac8a761665611e2e152eda75449753", "tag-1 tag-2 tag-3"},
	}

	checkAssertions(assertions, index, t)
}
