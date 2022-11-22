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
		argShId  string
		wantUrl  string
		wantIsOk bool
	}{
		{
			name:     "test#1",
			argShId:  "shortid",
			wantIsOk: true,
			wantUrl:  "https://Kirill.Znamenskiy.pw",
		},
		{
			name:     "test#2",
			argShId:  "aaaa",
			wantIsOk: false,
			wantUrl:  "",
		},
		{
			name:     "test#3",
			argShId:  "",
			wantIsOk: false,
			wantUrl:  "",
		},
	}
	stg := makeDefaultStorage(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUrl, gotIsOk := stg.Get(tt.argShId)
			assert.Equal(t, tt.wantIsOk, gotIsOk)
			gotUrlString := ""
			if gotUrl != nil {
				gotUrlString = gotUrl.String()
			}
			assert.Equal(t, tt.wantUrl, gotUrlString)
		})
	}
}

func TestInMemoryStorage_Put(t *testing.T) {
	tests := []struct {
		name    string
		argShId string
		argUrl  string
		want    bool
	}{
		{
			name:    "test#1",
			argShId: "aaa",
			argUrl:  "abfasb",
			want:    true,
		},
		{
			name:    "test#2",
			argShId: "bbb",
			argUrl:  "http://abfasb.org",
			want:    true,
		},
		{
			name:    "test#2",
			argShId: "ccc",
			argUrl:  "",
			want:    true,
		},
	}
	stg := makeDefaultStorage(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.argUrl)
			if err != nil {
				t.Fatal(err)
			}

			got := stg.Put(tt.argShId, u)
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
