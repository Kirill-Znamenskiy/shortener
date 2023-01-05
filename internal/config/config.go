package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Kirill-Znamenskiy/Shortener/internal/crypto"
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
	StorageFilePath string `env:"FILE_STORAGE_PATH" flag:"storage-file-path f" desc:"(env FILE_STORAGE_PATH) storage file path"`
	UserCookieName  string `env:"USER_COOKIE_NAME,default=kkk" flag:"user-cookie-name" desc:"(env USER_COOKIE_NAME) user cookie name"`
	SecretKey       string `env:"SECRET_KEY" flag:"secret-key" desc:"(env SECRET_KEY) server secret key"`
	DatabaseDsn     string `env:"DATABASE_DSN" flag:"database-dsn d" desc:"(env DATABASE_DSN) database dsn"`
	stg             storage.Storage
}

func (cfg *Config) SetStorage(stg storage.Storage) {
	cfg.stg = stg
}
func (cfg *Config) GetStorage() storage.Storage {
	if cfg.stg == nil {
		if cfg.DatabaseDsn != "" {
			cfg.stg = storage.NewDBStorage(context.TODO(), cfg.DatabaseDsn)
		} else if cfg.StorageFilePath != "" {
			cfg.stg = storage.NewFileStorage(cfg.StorageFilePath)
		} else {
			cfg.stg = storage.NewInMemoryStorage()
		}
	}
	return cfg.stg
}

func (cfg *Config) GetSecretKey() (ret []byte, err error) {
	storageSecretKey := cfg.GetStorage().GetSecretKey()
	if len(storageSecretKey) == 0 {
		if cfg.SecretKey == "" {
			ret, err = crypto.GenerateSecretKey(32)
			if err != nil {
				return nil, err
			}
		} else {
			ret = []byte(cfg.SecretKey)
		}
		err = cfg.GetStorage().PutSecretKey(ret)
		if err != nil {
			return nil, err
		}
		return
	} else {
		if cfg.SecretKey == "" || cfg.SecretKey == string(storageSecretKey) {
			ret = storageSecretKey
		} else {
			return nil, fmt.Errorf("its different secret keys in env and in storage")
		}
		return
	}
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
