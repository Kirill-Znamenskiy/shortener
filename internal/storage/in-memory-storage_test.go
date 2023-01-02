package storage

import (
	"github.com/Kirill-Znamenskiy/Shortener/pkg/kztests"
	"github.com/stretchr/testify/assert"
	"log"
	"net/url"
	"testing"
)

func makeDefaultStorage(t *testing.T) Storage {
	stg := NewInMemoryStorage()
	u, err := url.Parse("https://Kirill.Znamenskiy.pw")
	if err != nil {
		t.Fatal(err)
	}
	err = stg.Put("shortid", u)
	if err != nil {
		t.Fatal(err)
	}

	return stg
}

func TestInMemoryStorage_Get(t *testing.T) {
	stg := makeDefaultStorage(t)
	functions := stg.Get
	makeCheckResult1Func := func(expectedURL string) any {
		return func(t *testing.T, resURL *url.URL) bool {
			return resURL.String() == expectedURL
		}
	}
	testKits := []kztests.TestKit{
		{Arg: "shortid", Result1: makeCheckResult1Func("https://Kirill.Znamenskiy.pw"), Result2: true},
		{Arg: "aaa", Result1: assert.Nil, Result2: false},
		{Arg: "", Result1: assert.Nil, Result2: false},
	}
	kztests.RunTests(t, functions, testKits)
}

func TestInMemoryStorage_Put(t *testing.T) {
	stg := makeDefaultStorage(t)
	functions := stg.Put
	makeParseURLFunc := func(urlStr string) *url.URL {
		ret, err := url.Parse(urlStr)
		if err != nil {
			log.Fatal(err)
		}
		return ret
	}
	testKits := []kztests.TestKit{
		{Arg1: "aaa", Arg2: makeParseURLFunc("asdfas"), Result: assert.NoError},
		{Arg1: "bbb", Arg2: makeParseURLFunc("http://abfasb.org"), Result: assert.NoError},
		{Arg1: "ccc", Arg2: makeParseURLFunc(""), Result: assert.NoError},
		{Arg1: "ccc", Arg2: makeParseURLFunc("xxx"), Result: assert.Error},
	}
	kztests.RunTests(t, functions, testKits)
}

func TestNewInMemoryStorage(t *testing.T) {
	t.Run("Test 1", func(t *testing.T) {
		res := NewInMemoryStorage()
		assert.Empty(t, res.key2url)
	})
}
