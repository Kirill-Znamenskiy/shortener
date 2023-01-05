package storage_test

import (
	"context"
	"github.com/Kirill-Znamenskiy/Shortener/internal/config"
	"github.com/Kirill-Znamenskiy/Shortener/pkg/kztests"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage_Ping(t *testing.T) {
	cfg := new(config.Config)
	config.LoadFromEnv(context.TODO(), cfg)

	stg := cfg.GetStorage()

	functions := stg.Ping
	testKits := []kztests.TestKit{
		{Result1: assert.NoError},
	}
	kztests.RunTests(t, functions, testKits)
}
