package blogic

import (
	"github.com/Kirill-Znamenskiy/Shortener/internal/blogic/btypes"
	"github.com/Kirill-Znamenskiy/Shortener/internal/storage"
	"github.com/Kirill-Znamenskiy/Shortener/pkg/kztests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func makeDefaultStorage(t *testing.T) (stg storage.Storage) {
	stg, err := storage.NewInMemoryStorage()
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

func TestSaveNewURL(t *testing.T) {
	stg := makeDefaultStorage(t)
	sh := Shortener{stg: stg}
	functions := sh.SaveNewURL
	checkEmptyRecordKeyFunc := func(t *testing.T, rec *btypes.Record) bool {
		return rec == nil || rec.Key == ""
	}
	checkNonEmptyRecordKeyFunc := func(t *testing.T, rec *btypes.Record) bool {
		return rec.Key != ""
	}
	newUUID := uuid.New()
	user := btypes.User(&newUUID)
	testKits := []kztests.TestKit{
		{Arg1: user, Arg2: "https://Kirill.Znamenskiy.pw", Result1: checkNonEmptyRecordKeyFunc, Result2: assert.NoError},
		{Arg1: user, Arg2: "", Result1: checkNonEmptyRecordKeyFunc, Result2: assert.NoError},
		{Arg1: user, Arg2: "//yandex.ru", Result1: checkNonEmptyRecordKeyFunc, Result2: assert.NoError},
		{Arg1: user, Arg2: "://yandex.ru", Result1: checkEmptyRecordKeyFunc, Result2: assert.Error},
	}
	kztests.RunTests(t, functions, testKits)
}

func TestGetSavedURL(t *testing.T) {
	stg := makeDefaultStorage(t)
	sh := Shortener{stg: stg}
	functions := sh.GetSavedURL
	makeCheckResult1Func := func(expectedURL string) any {
		return func(t *testing.T, uuu *url.URL) bool {
			return uuu.String() == expectedURL
		}
	}
	testKits := []kztests.TestKit{
		{Arg: "shortid", Result1: makeCheckResult1Func("https://Kirill.Znamenskiy.pw"), Result2: true},
		{Arg: "aaa", Result1: assert.Nil, Result2: false},
		{Arg: "", Result1: assert.Nil, Result2: false},
	}
	kztests.RunTests(t, functions, testKits)
}
