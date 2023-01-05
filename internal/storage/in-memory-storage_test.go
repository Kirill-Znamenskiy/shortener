package storage

import (
	"github.com/Kirill-Znamenskiy/Shortener/internal/blogic/types"
	"github.com/Kirill-Znamenskiy/Shortener/pkg/kztests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"net/url"
	"testing"
)

func TestNewInMemoryStorage(t *testing.T) {
	t.Run("Test 1", func(t *testing.T) {
		res := NewInMemoryStorage()
		assert.True(t, res.IsEmpty())
	})
}

func makeDefaultStorage(t *testing.T) (stg Storage) {
	stg = NewInMemoryStorage()

	newUUID := uuid.New()
	user := types.User(&newUUID)

	u, err := url.Parse("https://Kirill.Znamenskiy.pw")
	if err != nil {
		t.Fatal(err)
	}
	err = stg.PutRecord(&types.Record{
		Key:         "shortid",
		OriginalURL: u,
		User:        user,
	})
	if err != nil {
		t.Fatal(err)
	}

	return
}

func TestInMemoryStorage_Get(t *testing.T) {
	stg := makeDefaultStorage(t)
	functions := stg.GetRecord
	makeCheckResult1Func := func(expectedURL string) any {
		return func(t *testing.T, rec *types.Record) bool {
			return rec.OriginalURL.String() == expectedURL
		}
	}
	testKits := []kztests.TestKit{
		{Arg: "shortid", Result1: makeCheckResult1Func("https://Kirill.Znamenskiy.pw")},
		{Arg: "aaa", Result1: assert.Nil},
		{Arg: "", Result1: assert.Nil},
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
	user := types.User(&newUUID)
	testKits := []kztests.TestKit{
		{Arg: &types.Record{Key: "aaa", OriginalURL: parseURLFunc("asdfas"), User: user}, Result: assert.NoError},
		{Arg: &types.Record{Key: "bbb", OriginalURL: parseURLFunc("http://abfasb.org"), User: user}, Result: assert.NoError},
		{Arg: &types.Record{Key: "ccc", OriginalURL: parseURLFunc(""), User: user}, Result: assert.NoError},
		{Arg: &types.Record{Key: "ccc", OriginalURL: parseURLFunc("xxx"), User: user}, Result: assert.Error},
	}
	kztests.RunTests(t, functions, testKits)
}
