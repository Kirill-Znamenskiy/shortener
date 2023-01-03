package blogic

import (
	"github.com/Kirill-Znamenskiy/Shortener/internal/storage"
	"github.com/Kirill-Znamenskiy/Shortener/pkg/kztests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func makeDefaultStorage(t *testing.T) (stg storage.Storage, userUUID *uuid.UUID) {
	stg = storage.NewInMemoryStorage()

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

func TestSaveNewURL(t *testing.T) {
	stg, userUUID := makeDefaultStorage(t)
	sh := Shortener{stg: stg}
	functions := sh.SaveNewURL
	testKits := []kztests.TestKit{
		{Arg1: userUUID, Arg2: "https://Kirill.Znamenskiy.pw", Result1: assert.NotEmpty, Result2: assert.NoError},
		{Arg1: userUUID, Arg2: "", Result1: assert.NotEmpty, Result2: assert.NoError},
		{Arg1: userUUID, Arg2: "//yandex.ru", Result1: assert.NotEmpty, Result2: assert.NoError},
		{Arg1: userUUID, Arg2: "://yandex.ru", Result1: "", Result2: assert.Error},
	}
	kztests.RunTests(t, functions, testKits)
}

func TestGetSavedURL(t *testing.T) {
	stg, userUUID := makeDefaultStorage(t)
	sh := Shortener{stg: stg}
	functions := sh.GetSavedURL
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
