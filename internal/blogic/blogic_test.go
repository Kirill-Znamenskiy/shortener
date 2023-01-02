package blogic

import (
	"github.com/Kirill-Znamenskiy/Shortener/internal/storage"
	"github.com/Kirill-Znamenskiy/Shortener/pkg/kztests"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func makeDefaultStorage(t *testing.T) storage.Storage {
	stg := storage.NewInMemoryStorage()
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

func TestSaveNewURL(t *testing.T) {
	functions := SaveNewURL
	stg := makeDefaultStorage(t)
	testKits := []kztests.TestKit{
		{Arg1: stg, Arg2: "https://Kirill.Znamenskiy.pw", Result1: assert.NotEmpty, Result2: assert.NoError},
		{Arg1: stg, Arg2: "", Result1: assert.NotEmpty, Result2: assert.NoError},
		{Arg1: stg, Arg2: "//yandex.ru", Result1: assert.NotEmpty, Result2: assert.NoError},
		{Arg1: stg, Arg2: "://yandex.ru", Result1: "", Result2: assert.Error},
	}
	kztests.RunTests(t, functions, testKits)
}

func TestGetSavedURL(t *testing.T) {
	functions := GetSavedURL
	makeCheckResult1Func := func(expectedURL string) any {
		return func(t *testing.T, resURL *url.URL) bool {
			return resURL.String() == expectedURL
		}
	}
	stg := makeDefaultStorage(t)
	testKits := []kztests.TestKit{
		{Arg1: stg, Arg2: "shortid", Result1: makeCheckResult1Func("https://Kirill.Znamenskiy.pw"), Result2: true},
		{Arg1: stg, Arg2: "aaa", Result1: assert.Nil, Result2: false},
		{Arg1: stg, Arg2: "", Result1: assert.Nil, Result2: false},
	}
	kztests.RunTests(t, functions, testKits)
}
