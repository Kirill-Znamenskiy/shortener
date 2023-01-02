package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Kirill-Znamenskiy/Shortener/internal/storage"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/sethvargo/go-envconfig"
	"github.com/spf13/pflag"
	"log"
	"os"
	"strings"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS,default=localhost:8080" flag:"server-address a" desc:"(env SERVER_ADDRESS) server address"`
	BaseURL         string `env:"BASE_URL,default=http://localhost:8080" flag:"base-url b" desc:"(env BASE_URL) base url"`
	StorageFilePath string `env:"FILE_STORAGE_PATH,default=localhost:8080" flag:"storage-file-path f" desc:"(env FILE_STORAGE_PATH) storage file path"`
	stg             storage.Storage
}

func (cfg *Config) SetStorage(stg storage.Storage) {
	cfg.stg = stg
}
func (cfg *Config) GetStorage() storage.Storage {
	if cfg.stg == nil {
		if cfg.StorageFilePath == "" {
			cfg.stg = storage.NewInMemoryStorage()
		} else {
			cfg.stg = storage.NewFileStorage(cfg.StorageFilePath)
		}
	}
	return cfg.stg
}

type myCustomEnvLookuper struct{}

func (*myCustomEnvLookuper) Lookup(key string) (string, bool) {
	val, ok := envconfig.OsLookuper().Lookup(key)

	// if env variable exists, but empty string - ignore it
	if ok && val == "" {
		ok = false
	}

	// if it consists only from empty spaces or/and quotes - ignore it too
	if ok {
		tmp := val
		tmp = strings.TrimSpace(tmp)
		tmp = strings.Trim(tmp, " '\"\t\n\v\f\r")
		tmp = strings.TrimSpace(tmp)
		if tmp == "" {
			ok = false
			val = ""
		}
	}

	return val, ok
}

func LoadFromEnv(ctx context.Context, cfg *Config) {
	err := envconfig.ProcessWith(ctx, cfg, new(myCustomEnvLookuper))
	if err != nil {
		log.Fatal(err)
	}
}

func ParseFlags(cfg *Config) {

	fs := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)

	err := gpflag.ParseTo(cfg, fs)
	if err != nil {
		log.Fatalln(err)
	}

	err = fs.Parse(os.Args[1:])
	if err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			fmt.Println("Help message finished.")
			os.Exit(0)
		}
		log.Fatalln(err)
	}
}

func ToPrettyString(cfg any) string {
	ret, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		log.Fatalln(err)
	}
	return string(ret)
}
