package storage

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func makeDefaultStorage(t *testing.T) Storage {
	stg := NewInMemoryStorage()
	u, err := url.Parse("https://Kirill.Znamenskiy.pw")
	if err != nil {
		t.Fatal(err)
	}
	stg.Put("shortid", u)

	return stg
}

func TestInMemoryStorage_Get(t *testing.T) {
	tests := []struct {
		name     string
		argShID  string
		wantURL  string
		wantIsOk bool
	}{
		{
			name:     "test#1",
			argShID:  "shortid",
			wantIsOk: true,
			wantURL:  "https://Kirill.Znamenskiy.pw",
		},
		{
			name:     "test#2",
			argShID:  "aaaa",
			wantIsOk: false,
			wantURL:  "",
		},
		{
			name:     "test#3",
			argShID:  "",
			wantIsOk: false,
			wantURL:  "",
		},
	}
	stg := makeDefaultStorage(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, gotIsOk := stg.Get(tt.argShID)
			assert.Equal(t, tt.wantIsOk, gotIsOk)
			gotURLString := ""
			if gotURL != nil {
				gotURLString = gotURL.String()
			}
			assert.Equal(t, tt.wantURL, gotURLString)
		})
	}
}

func TestInMemoryStorage_Put(t *testing.T) {
	tests := []struct {
		name    string
		argShID string
		argURL  string
		want    bool
	}{
		{
			name:    "test#1",
			argShID: "aaa",
			argURL:  "abfasb",
			want:    true,
		},
		{
			name:    "test#2",
			argShID: "bbb",
			argURL:  "http://abfasb.org",
			want:    true,
		},
		{
			name:    "test#2",
			argShID: "ccc",
			argURL:  "",
			want:    true,
		},
	}
	stg := makeDefaultStorage(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.argURL)
			if err != nil {
				t.Fatal(err)
			}

			got := stg.Put(tt.argShID, u)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewInMemoryStorage(t *testing.T) {
	t.Run("test#1", func(t *testing.T) {
		got := NewInMemoryStorage()
		assert.Empty(t, got.shid2url)
	})
}
