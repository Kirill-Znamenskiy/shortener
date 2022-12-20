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
	_, err = stg.Put("shortid", u)
	if err != nil {
		t.Fatal(err)
	}

	return stg
}

func TestInMemoryStorage_Get(t *testing.T) {
	tests := []struct {
		name     string
		argKey   string
		wantURL  string
		wantIsOk bool
	}{
		{
			name:     "test#1",
			argKey:   "shortid",
			wantIsOk: true,
			wantURL:  "https://Kirill.Znamenskiy.pw",
		},
		{
			name:     "test#2",
			argKey:   "aaaa",
			wantIsOk: false,
			wantURL:  "",
		},
		{
			name:     "test#3",
			argKey:   "",
			wantIsOk: false,
			wantURL:  "",
		},
	}
	stg := makeDefaultStorage(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, gotIsOk := stg.Get(tt.argKey)
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
		name   string
		argKey string
		argURL string
		want   bool
	}{
		{
			name:   "test#1",
			argKey: "aaa",
			argURL: "abfasb",
			want:   true,
		},
		{
			name:   "test#2",
			argKey: "bbb",
			argURL: "http://abfasb.org",
			want:   true,
		},
		{
			name:   "test#2",
			argKey: "ccc",
			argURL: "",
			want:   true,
		},
	}
	stg := makeDefaultStorage(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.argURL)
			if err != nil {
				t.Fatal(err)
			}

			got, err := stg.Put(tt.argKey, u)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewInMemoryStorage(t *testing.T) {
	t.Run("test#1", func(t *testing.T) {
		got := NewInMemoryStorage()
		assert.Empty(t, got.key2url)
	})
}
