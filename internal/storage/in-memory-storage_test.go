package storage

import (
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
		assert.Empty(t, res.GetSecretKey())
		assert.Empty(t, res.getSrcMap())
	})
}

func makeDefaultStorage(t *testing.T) (stg Storage, userUUID *uuid.UUID) {
	stg = NewInMemoryStorage()

	tmp := uuid.New()
	userUUID = &tmp

	u, err := url.Parse("https://Kirill.Znamenskiy.pw")
	if err != nil {
		t.Fatal(err)
	}
	err = stg.Put(userUUID, "shortid", u)
	if err != nil {
		t.Fatal(err)
	}

	return
}

func TestInMemoryStorage_Get(t *testing.T) {
	stg, userUUID := makeDefaultStorage(t)
	functions := stg.Get
	makeCheckResult1Func := func(expectedURL string) any {
		return func(t *testing.T, resURL *url.URL) bool {
			return resURL.String() == expectedURL
		}
	}
	testKits := []kztests.TestKit{
		{Arg1: userUUID, Arg2: "shortid", Result1: makeCheckResult1Func("https://Kirill.Znamenskiy.pw"), Result2: true},
		{Arg1: userUUID, Arg2: "aaa", Result1: assert.Nil, Result2: false},
		{Arg1: userUUID, Arg2: "", Result1: assert.Nil, Result2: false},
	}
	kztests.RunTests(t, functions, testKits)
}

func TestInMemoryStorage_Put(t *testing.T) {
	stg, userUUID := makeDefaultStorage(t)
	functions := stg.Put
	makeParseURLFunc := func(urlStr string) *url.URL {
		ret, err := url.Parse(urlStr)
		if err != nil {
			log.Fatal(err)
		}
		return ret
	}
	testKits := []kztests.TestKit{
		{Arg1: userUUID, Arg2: "aaa", Arg3: makeParseURLFunc("asdfas"), Result: assert.NoError},
		{Arg1: userUUID, Arg2: "bbb", Arg3: makeParseURLFunc("http://abfasb.org"), Result: assert.NoError},
		{Arg1: userUUID, Arg2: "ccc", Arg3: makeParseURLFunc(""), Result: assert.NoError},
		{Arg1: userUUID, Arg2: "ccc", Arg3: makeParseURLFunc("xxx"), Result: assert.Error},
	}
	kztests.RunTests(t, functions, testKits)
}
