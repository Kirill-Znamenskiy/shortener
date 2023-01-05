package storage

import (
	"github.com/Kirill-Znamenskiy/Shortener/internal/blogic/btypes"
	"github.com/Kirill-Znamenskiy/Shortener/pkg/kztests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/url"
	"testing"
)

func TestNewInMemoryStorage(t *testing.T) {
	t.Run("Test 1", func(t *testing.T) {
		res, err := NewInMemoryStorage()
		assert.NoError(t, err)
		assert.True(t, res.IsEmpty())
	})
}

func makeDefaultStorage(t *testing.T) (stg Storage) {
	stg, err := NewInMemoryStorage()
	require.NoError(t, err)

	newUUID := uuid.New()
	user := btypes.User(&newUUID)

	u, err := url.Parse("https://Kirill.Znamenskiy.pw")
	require.NoError(t, err)
	err = stg.PutRecord(&btypes.Record{
		Key:         "shortid",
		OriginalURL: u,
		User:        user,
	})
	require.NoError(t, err)

	return
}

func TestInMemoryStorage_Get(t *testing.T) {
	stg := makeDefaultStorage(t)
	functions := stg.GetRecord
	makeCheckResult1Func := func(expectedURL string) any {
		return func(t *testing.T, rec *btypes.Record) bool {
			return rec.OriginalURL.String() == expectedURL
		}
	}
	testKits := []kztests.TestKit{
		{Arg: "shortid", Result1: makeCheckResult1Func("https://Kirill.Znamenskiy.pw"), Result2: assert.NoError},
		{Arg: "aaa", Result1: assert.Nil, Result2: assert.NoError},
		{Arg: "", Result1: assert.Nil, Result2: assert.NoError},
	}
	kztests.RunTests(t, functions, testKits)
}

func TestInMemoryStorage_Put(t *testing.T) {
	stg := makeDefaultStorage(t)
	functions := stg.PutRecord
	parseURLFunc := func(urlStr string) *url.URL {
		ret, err := url.Parse(urlStr)
		if err != nil {
			log.Fatal(err)
		}
		return ret
	}
	newUUID := uuid.New()
	user := btypes.User(&newUUID)
	testKits := []kztests.TestKit{
		{Arg: &btypes.Record{Key: "aaa", OriginalURL: parseURLFunc("asdfas"), User: user}, Result: assert.NoError},
		{Arg: &btypes.Record{Key: "bbb", OriginalURL: parseURLFunc("http://abfasb.org"), User: user}, Result: assert.NoError},
		{Arg: &btypes.Record{Key: "ccc", OriginalURL: parseURLFunc(""), User: user}, Result: assert.NoError},
		{Arg: &btypes.Record{Key: "ccc", OriginalURL: parseURLFunc("xxx"), User: user}, Result: assert.Error},
	}
	kztests.RunTests(t, functions, testKits)
}
