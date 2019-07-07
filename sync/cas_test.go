package sync_test

import (
	"fmt"
	"path"
	"testing"

	"github.com/martinrue/s3-sync/sync"
)

func TestCreateCASKey(t *testing.T) {
	tests := []struct {
		Filename    string
		ExpectedKey string
	}{
		{"file-1.test", "2fa6e1156e6c6b2c808e9ef3c63f9c49ce809579200f2a40487949efd4febca4"},
		{"file-2.test", "99d5293841b29576ce680d53742c480f892e83ae441919b6f4f1f590879a5a24"},
		{"file-3.test", "23b2ff797a148f61f24c08cecb8182af89ac8a761665611e2e152eda75449753"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s", test.Filename), func(t *testing.T) {
			key, err := sync.CreateCASKey(path.Join("test_data/cas", test.Filename))
			if err != nil {
				t.Fatal(err)
			}

			if key != test.ExpectedKey {
				t.Fatalf("expected key %v, got %v", test.ExpectedKey, key)
			}
		})
	}
}
